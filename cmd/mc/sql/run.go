package sql

import (
	"database/sql"
	"fmt"
	"os"
	"reflect"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/sbchaos/opms/cmd/mc/internal"
	mcc "github.com/sbchaos/opms/external/mc"
	"github.com/sbchaos/opms/lib/cmdutil"
	"github.com/sbchaos/opms/lib/config"
	"github.com/sbchaos/opms/lib/table"
	"github.com/sbchaos/opms/lib/term"
)

type runSQL struct {
	cfg *config.Config

	query   string
	sqlFile string
}

func NewRunSQLCommand(cfg *config.Config) *cobra.Command {
	ec := &runSQL{cfg: cfg}

	cmd := &cobra.Command{
		Use:     "run",
		Short:   "Run a SQL query",
		Example: "opms mc sql run",
		RunE:    ec.RunE,
	}

	cmd.Flags().StringVarP(&ec.query, "query", "q", "", "Query to run")
	cmd.Flags().StringVarP(&ec.sqlFile, "file", "f", "", "Query filename to run")
	return cmd
}

func (r *runSQL) RunE(_ *cobra.Command, _ []string) error {
	t := term.FromEnv(0, 0)
	size, _ := t.Size(120)

	client, err := mcc.NewSQLClientFromConfig(r.cfg)
	if err != nil {
		return err
	}

	query := r.query
	if r.sqlFile != "" {
		bytes, err := cmdutil.ReadFile(r.sqlFile, os.Stdin)
		if err != nil {
			return err
		}
		query = string(bytes)
	}

	if query == "" {
		return fmt.Errorf("must specify query")
	}

	printer := table.New(os.Stdout, t.IsTerminalOutput(), size)
	err = runQuery(client, query, printer)
	if err != nil {
		return err
	}

	return printer.Render()
}

func runQuery(client *sql.DB, query string, printer table.Printer) error {
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
