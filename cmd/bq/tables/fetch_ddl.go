package tables

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"

	"cloud.google.com/go/bigquery"
	"github.com/spf13/cobra"
	"google.golang.org/api/iterator"

	"github.com/sbchaos/opms/external/gcp"
	"github.com/sbchaos/opms/lib/cmdutil"
	"github.com/sbchaos/opms/lib/config"
	"github.com/sbchaos/opms/lib/names"
	"github.com/sbchaos/opms/lib/printers/table"
	"github.com/sbchaos/opms/lib/term"
)

type fetchDDL struct {
	cfg       *config.Config
	name      string
	fileName  string
	outputDir string

	provider *gcp.ClientProvider
	printer  table.Printer
}

func NewFetchDDLCommand(cfg *config.Config) *cobra.Command {
	fetch := &fetchDDL{cfg: cfg}

	cmd := &cobra.Command{
		Use:     "fetch-ddl",
		Short:   "Fetch DDL for the table",
		Example: "opms bq tables count",
		RunE:    fetch.RunE,
	}

	cmd.Flags().StringVarP(&fetch.name, "name", "n", "", "Table name")
	cmd.Flags().StringVarP(&fetch.fileName, "filename", "f", "", "Filename with list of tables, - for stdin")
	cmd.Flags().StringVarP(&fetch.outputDir, "output-dir", "o", "", "Output directory")
	return cmd
}

func (m *fetchDDL) RunE(_ *cobra.Command, _ []string) error {
	client, err := gcp.NewClientProvider(m.cfg)
	if err != nil {
		return err
	}
	m.provider = client

	var tableNames []string
	if m.name == "" && m.fileName == "" {
		return errors.New("either --name or --filename is required")
	}

	if m.name != "" {
		tableNames = []string{m.name}
	}

	if m.fileName != "" {
		fields, err := cmdutil.ReadLines(m.fileName, os.Stdin)
		if err != nil {
			return err
		}

		tableNames = fields
	}
	t := term.FromEnv(0, 0)
	size, _ := t.Size(120)

	ctx, cancelFunc := context.WithTimeout(context.Background(), timeout)
	defer cancelFunc()

	printer := table.New(os.Stdout, t.IsTerminalOutput(), size)
	m.printer = printer

	printer.AddHeader([]string{"Table", "Status", "Error"})

	errs := make([]error, 0)
	for _, t1 := range tableNames {
		err = m.queryDDL(ctx, t1)
		if err != nil {
			errs = append(errs, err)
		}
	}

	printer.Render()
	if len(errs) != 0 {
		fmt.Println("Errors:")
		for _, err := range errs {
			fmt.Println("  " + err.Error())
		}
	}
	return nil
}

func (m *fetchDDL) queryDDL(ctx context.Context, tableName string) error {
	tb, err := names.FromTableName(tableName)
	if err != nil {
		return err
	}

	client, err := m.provider.GetClient(tb.Schema.ProjectID, driveScope)
	if err != nil {
		return err
	}

	query := fmt.Sprintf("select table_name, ddl from `%s.%s`.INFORMATION_SCHEMA.TABLES where table_name = '%s';", tb.Schema.ProjectID, tb.Schema.SchemaID, tb.TableID)
	q := client.Query(query)
	m.printer.AddField(tableName)

	it, err := q.Read(ctx)
	if err != nil {
		fetchFailure(m.printer, "permission issue")
		return fmt.Errorf("error while reading from bq: %w", err)
	}
	for {
		var row []bigquery.Value
		err = it.Next(&row)
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			fetchFailure(m.printer, "failure in query")
			return err
		}

		if len(row) != 2 {
			fetchFailure(m.printer, "too many rows")
		}
		content, ok := row[1].(string)
		if !ok {
			fetchFailure(m.printer, "unable to parse dll")
			return nil
		}

		toWritePath := tableName + ".sql"
		if m.outputDir != "" {
			toWritePath = path.Join(m.outputDir, toWritePath)
		}

		err = cmdutil.WriteFileAndDir(toWritePath, []byte(content))
		if err != nil {
			fetchFailure(m.printer, "failure in write file")
			return err
		}

		m.printer.AddField("success")
		m.printer.EndRow()
	}
	return nil
}

func fetchFailure(printer table.Printer, message string) {
	printer.AddField(message)
	printer.EndRow()
}
