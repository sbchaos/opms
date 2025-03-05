package tables

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/spf13/cobra"
	"google.golang.org/api/iterator"

	"github.com/sbchaos/opms/external/gcp"
	"github.com/sbchaos/opms/lib/cmdutil"
	"github.com/sbchaos/opms/lib/config"
	"github.com/sbchaos/opms/lib/names"
	"github.com/sbchaos/opms/lib/pool"
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
	workers     int
	mu          sync.Mutex
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
	cmd.Flags().IntVarP(&count.workers, "workers", "w", 1, "Number of parallel workers")

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

	tasks := make([]func() pool.JobResult[string], len(tableNames))
	for i, t1 := range tableNames {
		t1 := t1
		tasks[i] = func() pool.JobResult[string] {
			err = r.queryTable(ctx, provider, t1, printer)
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

func (r *countCommand) queryTable(ctx context.Context, provider *gcp.ClientProvider, tableName string, printer table.Printer) error {
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

	it, err := q.Read(ctx)
	if err != nil {
		r.mu.Lock()
		printer.AddField(tableName)
		addFailure(printer, "permission issue")
		r.mu.Unlock()
		return fmt.Errorf("error while reading from bq: %w", err)
	}
	for {
		var row []bigquery.Value
		err = it.Next(&row)
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			r.mu.Lock()
			printer.AddField(tableName)
			addFailure(printer, "failure in query")
			r.mu.Unlock()
			return err
		}

		if len(row) != 1 {
			r.mu.Lock()
			printer.AddField(tableName)
			addFailure(printer, "too many rows")
			r.mu.Unlock()
		}

		r.mu.Lock()
		printer.AddField(tableName)
		printer.AddField(fmt.Sprintf("%v", row[0]))
		printer.AddField("")
		printer.EndRow()
		r.mu.Unlock()
	}
	return nil
}

func addFailure(printer table.Printer, msg string) {
	printer.AddField("-1")
	printer.AddField(msg)
	printer.EndRow()
}
