package tables

import (
	"database/sql"
	"os"
	"reflect"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/sbchaos/opms/cmd/mc/internal"
	mcc "github.com/sbchaos/opms/external/mc"
	"github.com/sbchaos/opms/lib/config"
	"github.com/sbchaos/opms/lib/table"
	"github.com/sbchaos/opms/lib/term"
)

var (
	deadlineErr = "context deadline exceeded"
	queryFields = `SELECT * FROM `
)

type readTable struct {
	cfg *config.Config

	name string
}

// NewReadTableCommand checks if the tables exist
func NewReadTableCommand(cfg *config.Config) *cobra.Command {
	ec := &readTable{cfg: cfg}

	cmd := &cobra.Command{
		Use:     "read",
		Short:   "Read data from a table",
		Example: "opms mc table read",
		RunE:    ec.RunE,
	}

	cmd.Flags().StringVarP(&ec.name, "name", "n", "", "Table name")
	cmd.MarkFlagRequired("name")
	return cmd
}

func (r *readTable) RunE(_ *cobra.Command, _ []string) error {
	t := term.FromEnv(0, 0)
	size, _ := t.Size(120)

	client, err := mcc.NewSQLClientFromConfig(r.cfg)
	if err != nil {
		return err
	}

	printer := table.New(os.Stdout, t.IsTerminalOutput(), size)
	err = runQuery(client, r.name, printer)
	if err != nil {
		return err
	}

	return printer.Render()
}

func runQuery(client *sql.DB, name string, printer table.Printer) error {
	query := queryFields + name + ";"
	rows, err := client.Query(query)
	if err != nil {
		return err
	}

	cols, err := rows.Columns()
	if err != nil {
		return err
	}

	headers := append([]string{"Row"}, cols...)
	printer.AddHeader(headers)

	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		return err
	}

	record := make([]interface{}, len(columnTypes))
	for i, columnType := range columnTypes {
		record[i] = reflect.New(columnType.ScanType()).Interface()
	}

	rowNum := 1
	for rows.Next() {
		err = rows.Scan(record...)
		if err != nil {
			return err
		}

		printer.AddField(strconv.Itoa(rowNum))
		for _, r := range record {
			str := internal.ToString(r)
			printer.AddField(str)
		}
		printer.EndRow()
		rowNum++
	}

	return nil
}
