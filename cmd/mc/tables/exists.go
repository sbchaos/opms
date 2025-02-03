package tables

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/aliyun/aliyun-odps-go-sdk/odps"
	"github.com/spf13/cobra"

	mcc "github.com/sbchaos/opms/external/mc"
	"github.com/sbchaos/opms/lib/cmdutil"
	"github.com/sbchaos/opms/lib/config"
	"github.com/sbchaos/opms/lib/list"
	"github.com/sbchaos/opms/lib/table"
	"github.com/sbchaos/opms/lib/term"
)

var (
	connErr            = "connection reset by peer"
	notExistErrorRegex = regexp.MustCompile(`Table (\S*) does not exist`)
)

type existsCommand struct {
	cfg *config.Config

	name     string
	fileName string
}

// NewExistsCommand checks if the tables exist
func NewExistsCommand(cfg *config.Config) *cobra.Command {
	ec := &existsCommand{cfg: cfg}

	cmd := &cobra.Command{
		Use:     "exists",
		Short:   "Show if the table exists in maxcompute",
		Example: "opms mc tables exists",
		RunE:    ec.RunE,
	}

	cmd.Flags().StringVarP(&ec.name, "name", "n", "", "Table name")
	cmd.Flags().StringVarP(&ec.fileName, "filename", "f", "", "Filename with list of tables, - for stdin")
	return cmd
}

func (r *existsCommand) RunE(_ *cobra.Command, _ []string) error {
	t := term.FromEnv(0, 0)
	size, _, err := t.Size()
	if err != nil {
		size = 120
	}

	client, err := mcc.NewClientFromConfig(r.cfg)
	if err != nil {
		return err
	}

	mapping := make(map[mcc.ProjectSchema][]string)
	if r.name != "" {
		ps, name, err := mcc.SplitParts(r.name)
		if err != nil {
			return err
		}

		mapping[ps] = []string{name}

		if r.fileName != "" {
			return errors.New("--filename flag cannot be used along with name")
		}
	}

	if r.fileName != "" {
		content, err := cmdutil.ReadFile(r.fileName, os.Stdin)
		if err != nil {
			return err
		}
		fields := strings.Fields(string(content))
		for _, field := range fields {
			ps, name, err := mcc.SplitParts(field)
			if err != nil {
				fmt.Printf("ignoring invalid table name %s\n", field)
			}

			mapping[ps] = append(mapping[ps], name)
		}
	}

	printer := table.New(os.Stdout, t.IsTerminalOutput(), size)

	for ps, tables := range mapping {
		client.SetDefaultProjectName(ps.ProjectName)
		client.SetCurrentSchemaName(ps.SchemaName)

		lenTables := len(tables)
		for i := 0; i < lenTables; i = i + 100 {
			end := min(i+100, lenTables)
			err := forN(100, printer, client.Tables(), ps, tables[i:end])
			if err != nil {
				return err
			}
		}
	}

	err = printer.Render()
	if err != nil {
		return fmt.Errorf("failed to print table: %w", err)
	}
	return nil
}

func forN(step int, printer table.Printer, tabs *odps.Tables, ps mcc.ProjectSchema, tables []string) error {
	lenTables := len(tables)
	for i := 0; i < lenTables; i = i + step {
		end := min(i+step, lenTables)
		batch := tables[i:end]

		for {
			loadTables, err := tabs.BatchLoadTables(batch)
			if err != nil {
				if strings.Contains(err.Error(), connErr) {
					time.Sleep(500 * time.Millisecond)
					continue
				}

				submatch := notExistErrorRegex.FindStringSubmatch(err.Error())
				if len(submatch[0]) > 0 {
					failureStatus(printer, ps, submatch[1])
					batch = list.Remove(batch, submatch[1])
					continue
				}
				return err
			} else {
				successStatus(printer, ps, loadTables)
				break
			}
		}
	}
	return nil
}

func failureStatus(printer table.Printer, ps mcc.ProjectSchema, table string) {
	printer.AddField(" ❌ ")
	printer.AddField(ps.Table(table))
	printer.EndRow()
}

func successStatus(printer table.Printer, ps mcc.ProjectSchema, tables []*odps.Table) {
	printer.AddHeader([]string{"Exists", "Table Name"})

	for _, t1 := range tables {
		printer.AddField(" ✅ ")
		printer.AddField(ps.Table(t1.Name()))
		printer.EndRow()
	}
}
