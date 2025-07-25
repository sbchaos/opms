package external_table

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/spf13/cobra"

	"github.com/sbchaos/opms/cmd/optimus/internal/parse"
	"github.com/sbchaos/opms/cmd/optimus/internal/resource"
	"github.com/sbchaos/opms/external/gcp"
	"github.com/sbchaos/opms/external/gsheet"
	"github.com/sbchaos/opms/lib/cmdutil"
	"github.com/sbchaos/opms/lib/config"
	"github.com/sbchaos/opms/lib/names"
	"github.com/sbchaos/opms/lib/printers/table"
	"github.com/sbchaos/opms/lib/term"
)

var (
	Success = " ✅ "
	Failed  = " ❌ "

	ErrNotReq = errors.New("not required")
)

//go:embed resource.yaml.tmpl
var yamlSpec []byte

type etCommand struct {
	cfg *config.Config

	fileName string
	dirName  string

	typeMapJson string
	typeMap     map[string]string

	projMapJson string
	projMap     map[string]string

	required     string
	requiredList map[string]string

	fixSheetRange bool
	provider      *gcp.ClientProvider

	tmpl   *template.Template
	parser *parse.DDLParser
}

func NewExternalTableCommand(cfg *config.Config) *cobra.Command {
	ec := &etCommand{cfg: cfg}

	cmd := &cobra.Command{
		Use:     "external-table",
		Short:   "Generate external table spec from bq query/queries",
		Example: "opms opt generate external-table --query /path/to/file.sql",
		PreRunE: ec.PreRunE,
		RunE:    ec.RunE,
	}

	cmd.Flags().StringVarP(&ec.fileName, "query", "q", "", "Query file name, - for stdin")
	cmd.Flags().StringVarP(&ec.dirName, "query-dir", "d", "", "Directory with queries")

	cmd.Flags().StringVarP(&ec.typeMapJson, "type-map", "t", "", "Mapping json of BQ to maxcompute type")
	cmd.Flags().StringVarP(&ec.projMapJson, "proj-map", "p", "", "Mapping json of BQ to maxcompute projects")
	cmd.Flags().StringVarP(&ec.required, "required", "r", "", "List of required tables, - for stdin")
	cmd.Flags().BoolVarP(&ec.fixSheetRange, "fix-range", "f", false, "Try to fix sheet range")
	cmd.MarkFlagRequired("type-map")
	cmd.MarkFlagRequired("proj-map")
	return cmd
}

func (r *etCommand) PreRunE(_ *cobra.Command, _ []string) error {
	var err error
	tmpl, err := template.New("maxcompute_migrated_spec").Parse(string(yamlSpec))
	if err != nil {
		return err
	}
	r.tmpl = tmpl

	if r.fileName == "" && r.dirName == "" {
		return fmt.Errorf("must provide either a file or a directory")
	}

	if r.fixSheetRange {
		p1, errProvider := gcp.NewClientProvider(r.cfg)
		if errProvider != nil {
			return errProvider
		}
		r.provider = p1
	}

	r.typeMap = map[string]string{}
	if r.typeMapJson != "" {
		mapping := make(map[string]string)
		err2 := cmdutil.ReadJsonFile(r.typeMapJson, os.Stdin, mapping)
		if err2 != nil {
			return err2
		}
		r.typeMap = mapping
	} else {
		fmt.Println("No mapping found for table types")
	}

	r.projMap = map[string]string{}
	if r.projMapJson != "" {
		mapping := make(map[string]string)
		err2 := cmdutil.ReadJsonFile(r.projMapJson, os.Stdin, mapping)
		if err2 != nil {
			return err2
		}
		r.projMap = mapping
	} else {
		fmt.Println("No mapping found for project")
	}

	reqFn := func(_ string) error {
		return nil
	}

	r.requiredList = map[string]string{}
	if r.required != "" {
		content, err := cmdutil.ReadFile(r.required, os.Stdin)
		if err != nil {
			return err
		}
		for _, st := range strings.Fields(string(content)) {
			r.requiredList[st] = st
		}

		reqFn = func(name string) error {
			if r.requiredList[name] != "" {
				return nil
			}
			return ErrNotReq
		}
	} else {
		fmt.Println("Processing all queries")
	}

	r.parser, err = parse.NewDDLParser(nil, reqFn)
	if err != nil {
		return err
	}

	return err
}

func (r *etCommand) RunE(_ *cobra.Command, _ []string) error {
	t := term.FromEnv(0, 0)
	size, _ := t.Size(120)
	printer := table.New(os.Stdout, t.IsTerminalOutput(), size)
	printer.AddHeader([]string{"Status", "Name", "Error"})

	var errs []error
	if r.fileName != "" {
		content, err := cmdutil.ReadFile(r.fileName, os.Stdin)
		if err != nil {
			return err
		}
		query := string(content)
		err = r.processQuery(r.fileName, query, printer)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if r.dirName != "" {
		errDirs := r.processDirectory(printer)
		errs = errDirs
	}

	printer.Render()
	if len(errs) > 0 {
		fmt.Println("Errors:")
		for _, e := range errs {
			fmt.Println(e)
		}
	}
	return nil
}

func (r *etCommand) processDirectory(printer table.Printer) []error {
	entries, err := os.ReadDir(r.dirName)
	if err != nil {
		addFailureRow(printer, r.dirName, "error reading dir")
		return nil
	}

	errs := make([]error, 0)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filePath := path.Join(r.dirName, entry.Name())
		content, err := os.ReadFile(filePath)
		if err != nil {
			addFailureRow(printer, filePath, "error opening file")
			continue
		}

		err = r.processQuery(entry.Name(), string(content), printer)
		if err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

func (r *etCommand) processQuery(name, query string, printer table.Printer) error {
	bqET, err := r.parser.ParseExternalTable(query)
	if err != nil {
		toName := name
		msg := "error parsing query"

		if bqET != nil && bqET.FullName() != "" {
			toName = bqET.FullName()
		}

		if err == ErrNotReq {
			msg = "not required"
		}

		addFailureRow(printer, toName, msg)
		return err
	}
	tableName := bqET.FullName()

	y1, err := resource.MapExternalTable(bqET, r.projMap, r.typeMap)
	if err != nil {
		addFailureRow(printer, tableName, "error converting to maxcompute type")
		return err
	}

	if r.fixSheetRange {
		r.fixRangeIfMissing(y1)
	}

	proj := y1.Et.Project
	projVar, ok := r.projMap[proj+"_multi"]
	if ok {
		y1.Et.Project = projVar
	}

	yctx := &YamlContext{
		Et:      y1.Et,
		Table:   y1.Table,
		OldName: y1.OldName,
		Labels:  nil,
	}

	err = r.WriteResource(yctx)
	if err != nil {
		addFailureRow(printer, tableName, "error writing file")
		return err
	}
	printer.AddField(Success)
	printer.AddField(tableName)
	printer.AddField("")
	printer.EndRow()
	return nil
}

func (r *etCommand) fixRangeIfMissing(met *resource.MappedExtTable) {
	sheetURI := met.Et.Source.SourceURIs[0]
	rng := met.Et.Source.Range
	if rng == "" || (strings.Contains(rng, ":") && !strings.Contains(rng, "!")) {
		proj, _, found := strings.Cut(met.OldName, ".")
		if !found {
			return
		}
		service, err2 := r.provider.GetSheetsClient(proj)
		if err2 != nil {
			fmt.Println(err2)
			return
		}
		name, err := gsheet.GetSheetName(service, sheetURI)
		if err == nil {
			if rng != "" {
				name += "!" + rng
			}
			met.Et.Source.Range = name
		} else {
			fmt.Println("Error getting sheet name:", err)
		}
	}
}

func addFailureRow(printer table.Printer, name string, msg string) {
	printer.AddField(Failed)
	printer.AddField(name)
	printer.AddField(msg)
	printer.EndRow()
}

func (r *etCommand) WriteResource(ymCtx *YamlContext) error {
	content, err := GetContent(r.tmpl, ymCtx)
	if err != nil {
		return err
	}
	filePath := path.Join("generated", "maxcompute", ymCtx.Table.String(), "resource.yaml")
	return cmdutil.WriteFileAndDir(filePath, []byte(content))
}

func GetContent(tmpl *template.Template, ctx *YamlContext) (string, error) {
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, ctx); err != nil {
		msg := fmt.Sprintf("unable to compile template for resource %s : %s", ctx.Et.Name, err.Error())
		return "", errors.New(msg)
	}

	return buf.String(), nil
}

type YamlContext struct {
	Et      *resource.ExternalTable
	Table   names.Table
	OldName string
	Labels  map[string]string
}
