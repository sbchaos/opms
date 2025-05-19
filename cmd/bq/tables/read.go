package tables

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"

	"cloud.google.com/go/bigquery"
	"github.com/spf13/cobra"
	"google.golang.org/api/iterator"

	"github.com/sbchaos/opms/external/gcp"
	"github.com/sbchaos/opms/lib/config"
	"github.com/sbchaos/opms/lib/names"
	"github.com/sbchaos/opms/lib/printers/table"
	"github.com/sbchaos/opms/lib/term"
)

type readCommand struct {
	cfg *config.Config

	name string
}

// NewReadCommand initializes command to read number of rows in table
func NewReadCommand(cfg *config.Config) *cobra.Command {
	read := &readCommand{cfg: cfg}

	cmd := &cobra.Command{
		Use:     "read",
		Short:   "Read the rows in table",
		Example: "opms bq tables read",
		RunE:    read.RunE,
	}

	cmd.Flags().StringVarP(&read.name, "name", "n", "", "Table name")
	cmd.MarkFlagRequired("name")
	return cmd
}

func (r *readCommand) RunE(_ *cobra.Command, _ []string) error {
	provider, err := gcp.NewClientProvider(r.cfg)
	if err != nil {
		return err
	}

	t := term.FromEnv(0, 0)
	size, _ := t.Size(120)

	ctx, cancelFunc := context.WithTimeout(context.Background(), timeout)
	defer cancelFunc()

	printer := table.New(os.Stdout, t.IsTerminalOutput(), size)

	err = readTable(ctx, provider, r.name, printer)
	if err != nil {
		return err
	}

	printer.Render()
	return nil
}

func readTable(ctx context.Context, provider *gcp.ClientProvider, tableName string, printer table.Printer) error {
	tb, err := names.FromTableName(tableName)
	if err != nil {
		return err
	}

	client, err := provider.GetClient(tb.Schema.ProjectID, driveScope)
	if err != nil {
		return err
	}

	qr := `SELECT * FROM ` + tableName
	q := client.Query(qr)

	it, err := q.Read(ctx)
	if err != nil {
		return fmt.Errorf("error while reading from bq: %w", err)
	}

	headers := []string{"row"}
	for _, field := range it.Schema {
		headers = append(headers, field.Name)
	}

	rowNum := 1
	for {
		var row []bigquery.Value
		err = it.Next(&row)
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return err
		}
		printer.AddField(strconv.Itoa(rowNum))
		decodeRow(row, printer)
		rowNum++
	}
	return nil
}

func decodeRow(row []bigquery.Value, printer table.Printer) {
	for _, field := range row {
		val := fmt.Sprintf("%v", field)
		printer.AddField(val)
	}
	printer.EndRow()
}
