package verify

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/spf13/cobra"

	mcc "github.com/sbchaos/opms/external/mc"
	"github.com/sbchaos/opms/lib/cmdutil"
	"github.com/sbchaos/opms/lib/config"
	"github.com/sbchaos/opms/lib/pool"
	"github.com/sbchaos/opms/lib/table"
	"github.com/sbchaos/opms/lib/term"
)

var (
	maxRecordLimit = int64(600000)

	countSQL    = `SELECT count(*) FROM `
	queryFields = `SELECT * FROM `
)

type externalTableCommand struct {
	cfg *config.Config

	mu      *sync.Mutex
	workers int

	name     string
	fileName string
}

// NewExternalTableCommand checks if the tables exist
func NewExternalTableCommand(cfg *config.Config) *cobra.Command {
	ec := &externalTableCommand{cfg: cfg}

	cmd := &cobra.Command{
		Use:     "externalTable",
		Short:   "Verify the externalTable in maxcompute",
		Example: "opms mc verify externalTable",
		RunE:    ec.RunE,
	}

	cmd.Flags().StringVarP(&ec.name, "name", "n", "", "Table name")
	cmd.Flags().StringVarP(&ec.fileName, "filename", "f", "", "Filename with list of tables, - for stdin")
	cmd.Flags().IntVarP(&ec.workers, "workers", "w", 1, "Number of parallel workers")
	return cmd
}

func (r *externalTableCommand) RunE(_ *cobra.Command, _ []string) error {
	r.mu = &sync.Mutex{}
	t := term.FromEnv(0, 0)
	size, _ := t.Size(120)

	client, err := mcc.NewSQLClientFromConfig(r.cfg)
	if err != nil {
		return err
	}

	tables := make([]string, 0)
	if r.name != "" {
		tables = append(tables, r.name)
		if r.fileName != "" {
			return errors.New("--filename flag cannot be used along with name")
		}
	}

	if r.fileName != "" {
		fields, err := cmdutil.ReadLines(r.fileName, os.Stdin)
		if err != nil {
			return err
		}

		tables = append(tables, fields...)
	}

	printer := table.New(os.Stdout, t.IsTerminalOutput(), size)
	printer.AddHeader([]string{"Status", "Table Name", "COUNT*", "SIZE", "Error"})

	tasks := make([]func() pool.JobResult[string], len(tables))
	for i, t1 := range tables {
		t1 := t1
		tasks[i] = func() pool.JobResult[string] {
			err = r.Validate(client, printer, t1)
			return pool.JobResult[string]{
				Output: t1,
				Err:    err,
			}
		}
	}

	outchan := pool.RunWithWorkers(r.workers, tasks)

	err = printer.Render()
	if err != nil {
		return fmt.Errorf("failed to print table: %w", err)
	}

	fmt.Println()
	fmt.Println()
	fmt.Fprintln(os.Stderr, "Error(s) encountered:")
	for out := range outchan {
		if out.Err != nil {
			fmt.Fprintf(os.Stderr, "Name: %s, Err: %s\n", out.Output, out.Err)
		}
	}

	return nil
}

func (r *externalTableCommand) Validate(client *sql.DB, printer table.Printer, name string) error {
	countStar := int64(-1)
	countRow := int64(0)

	res, err := runCountStar(client, name)
	if err == nil {
		countStar = res
	}

	if countStar > 0 && countStar < maxRecordLimit {
		res2, err := runCount(client, name)
		if err == nil {
			countRow = res2
		}
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	if countStar == countRow {
		printer.AddField(" ✅ ")
	} else if countStar > maxRecordLimit {
		printer.AddField(" ❗ ")
	} else {
		printer.AddField(" ❌ ")
	}
	printer.AddField(name)
	printer.AddField(strconv.FormatInt(countStar, 10))
	printer.AddField(strconv.FormatInt(countRow, 10))
	if err != nil {
		errFirstPart := err.Error()[0:20]
		printer.AddField(errFirstPart)
	}
	printer.EndRow()
	return err
}

func runCountStar(client *sql.DB, name string) (int64, error) {
	query := countSQL + name + ";"
	rows, err := client.Query(query)
	if err != nil {
		return -1, err
	}

	rowCount := int64(0)
	for rows.Next() {
		err = rows.Scan(&rowCount)
		if err != nil {
			return -1, err
		}
	}

	return rowCount, nil
}

func runCount(client *sql.DB, name string) (int64, error) {
	query := queryFields + name + ";"
	rows, err := client.Query(query)
	if err != nil {
		return -1, err
	}

	rowCount := int64(0)
	for rows.Next() {
		rowCount++
	}

	return rowCount, nil
}
