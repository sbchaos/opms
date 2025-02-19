package tables

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/spf13/cobra"
	"google.golang.org/api/iterator"

	"github.com/sbchaos/opms/external/bq"
	"github.com/sbchaos/opms/lib/cmdutil"
	"github.com/sbchaos/opms/lib/config"
	"github.com/sbchaos/opms/lib/table"
	"github.com/sbchaos/opms/lib/term"
)

var timeout = time.Minute * 10

type countCommand struct {
	cfg *config.Config

	name     string
	fileName string
}

// NewCountCommand initializes command to count number of rows in table
func NewCountCommand(cfg *config.Config) *cobra.Command {
	count := &countCommand{cfg: cfg}

	cmd := &cobra.Command{
		Use:     "count",
		Short:   "Count the tables",
		Example: "opms bq tables count",
		RunE:    count.RunE,
	}

	cmd.Flags().StringVarP(&count.name, "name", "n", "", "Table name")
	cmd.Flags().StringVarP(&count.fileName, "filename", "f", "", "Filename with list of tables, - for stdin")

	return cmd
}

func (r *countCommand) RunE(_ *cobra.Command, _ []string) error {
	client, err := bq.NewClientFromConfig(r.cfg)
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

	if r.fileName != "" {
		content, err := cmdutil.ReadFile(r.fileName, os.Stdin)
		if err != nil {
			return err
		}

		fields := strings.Fields(string(content))
		for _, field := range fields {
			tableNames = append(tableNames, field)
		}
	}
	t := term.FromEnv(0, 0)
	size, _ := t.Size(120)

	ctx, cancelFunc := context.WithTimeout(context.Background(), timeout)
	defer cancelFunc()

	printer := table.New(os.Stdout, t.IsTerminalOutput(), size)
	printer.AddHeader([]string{"Table", "Count", "Error"})

	for _, t1 := range tableNames {
		err = queryTable(ctx, client, t1, printer)
		if err != nil {
			fmt.Println(err)
		}
	}

	err = printer.Render()
	if err != nil {
		return fmt.Errorf("failed to print table: %w", err)
	}
	return nil
}

func queryTable(ctx context.Context, client *bq.Client, tableName string, printer table.Printer) error {
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
