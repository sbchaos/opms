package tables

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/spf13/cobra"
	"google.golang.org/api/iterator"

	"github.com/sbchaos/opms/external/gcp"
	"github.com/sbchaos/opms/lib/cmdutil"
	"github.com/sbchaos/opms/lib/config"
	"github.com/sbchaos/opms/lib/names"
	"github.com/sbchaos/opms/lib/table"
	"github.com/sbchaos/opms/lib/term"
)

var timeout = time.Minute * 10

var driveScope = "https://www.googleapis.com/auth/drive"

type countCommand struct {
	cfg *config.Config

	name     string
	fileName string

	mappingJson string
}

// NewCountCommand initializes command to count number of rows in table
func NewCountCommand(cfg *config.Config) *cobra.Command {
	count := &countCommand{cfg: cfg}

	cmd := &cobra.Command{
		Use:     "count",
		Short:   "Count the rows in table",
		Example: "opms bq tables count",
		RunE:    count.RunE,
	}

	cmd.Flags().StringVarP(&count.name, "name", "n", "", "Table name")
	cmd.Flags().StringVarP(&count.fileName, "filename", "f", "", "Filename with list of tables, - for stdin")
	cmd.Flags().StringVarP(&count.mappingJson, "mapping", "m", "", "Project mapping for the names")

	return cmd
}

func (r *countCommand) RunE(_ *cobra.Command, _ []string) error {
	provider, err := gcp.NewClientProvider(r.cfg)
	if err != nil {
		return err
	}

	var tableNames []string
	if r.name == "" && r.fileName == "" {
		return errors.New("either --name or --filename is required")
	}

	if r.name != "" {
		tableNames = []string{r.name}
	}

	projectMapping := map[string]string{}
	if r.mappingJson != "" {
		err = cmdutil.ReadJsonFile(r.mappingJson, os.Stdin, &projectMapping)
		if err != nil {
			return err
		}
	}

	if r.fileName != "" {
		fields, err := cmdutil.ReadLines(r.fileName, os.Stdin)
		if err != nil {
			return err
		}

		mapNames, err := names.MapNames(projectMapping, nil, fields)
		if err != nil {
			return err
		}
		tableNames = mapNames
	}
	t := term.FromEnv(0, 0)
	size, _ := t.Size(120)

	ctx, cancelFunc := context.WithTimeout(context.Background(), timeout)
	defer cancelFunc()

	printer := table.New(os.Stdout, t.IsTerminalOutput(), size)
	printer.AddHeader([]string{"Table", "Count", "Error"})

	errs := make([]error, 0)
	for _, t1 := range tableNames {
		err = queryTable(ctx, provider, t1, printer)
		if err != nil {
			errs = append(errs, err)
		}
	}

	err = printer.Render()
	if len(errs) != 0 {
		fmt.Println("Errors:")
		for _, err := range errs {
			fmt.Println("  " + err.Error())
		}
	}
	return nil
}

func queryTable(ctx context.Context, provider *gcp.ClientProvider, tableName string, printer table.Printer) error {
	tb, err := names.FromTableName(tableName)
	if err != nil {
		return err
	}

	client, err := provider.GetClient(tb.Schema.ProjectID, driveScope)
	if err != nil {
		return err
	}

	qr := `SELECT COUNT(*) FROM ` + tableName
	q := client.Query(qr)
	printer.AddField(tableName)

	it, err := q.Read(ctx)
	if err != nil {
		addFailure(printer, "permission issue")
		return fmt.Errorf("error while reading from bq: %w", err)
	}
	for {
		var row []bigquery.Value
		err = it.Next(&row)
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			addFailure(printer, "failure in query")
			return err
		}

		if len(row) != 1 {
			addFailure(printer, "too many rows")
		}

		printer.AddField(fmt.Sprintf("%v", row[0]))
		printer.AddField("")
		printer.EndRow()
	}
	return nil
}

func addFailure(printer table.Printer, msg string) {
	printer.AddField("-1")
	printer.AddField(msg)
	printer.EndRow()
}
