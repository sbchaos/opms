package verify

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/aliyun/aliyun-odps-go-sdk/odps"
	"github.com/spf13/cobra"

	mcc "github.com/sbchaos/opms/external/mc"
	"github.com/sbchaos/opms/lib/cmdutil"
	"github.com/sbchaos/opms/lib/config"
	"github.com/sbchaos/opms/lib/pool"
	"github.com/sbchaos/opms/lib/table"
	"github.com/sbchaos/opms/lib/term"
)

var (
	connErr     = "connection reset by peer"
	deadlineErr = "context deadline exceeded"

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
	size, _, err := t.Size()
	if err != nil {
		size = 120
	}

	client, err := mcc.NewClientFromConfig(r.cfg)
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
		content, err := cmdutil.ReadFile(r.fileName, os.Stdin)
		if err != nil {
			return err
		}

		fields := strings.Fields(string(content))
		for _, field := range fields {
			tables = append(tables, field)
		}
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

func (r *externalTableCommand) Validate(client *odps.Odps, printer table.Printer, name string) error {
	countStar := -1
	countRow := 0

	countStarQry := countSQL + name + ";"
	res, err := runQuery(client, countStarQry)
	if err == nil {
		parts := strings.Split(res, "\n")
		if len(parts) > 1 {
			val, err2 := strconv.Atoi(parts[1])
			if err2 == nil {
				countStar = val
			}
		}
	}

	countFieldsQry := queryFields + name + ";"
	res2, err := runQuery(client, countFieldsQry)
	if err == nil {
		parts := strings.Split(res2, "\n")
		for i, field := range parts {
			if field != "" && i > 0 {
				countRow++
			}
		}

	}

	r.mu.Lock()
	defer r.mu.Unlock()
	if countStar == countRow {
		printer.AddField(" ✅ ")
	} else if countRow == 10000 {
		printer.AddField(" ❗️ ")
	} else {
		printer.AddField(" ❌ ")
	}
	printer.AddField(name)
	printer.AddField(strconv.Itoa(countStar))
	printer.AddField(strconv.Itoa(countRow))
	if err != nil {
		errFirstPart := err.Error()[0:20]
		printer.AddField(errFirstPart)
	}
	printer.EndRow()
	return err
}

func runQuery(client *odps.Odps, query string) (string, error) {
	i := 5
	for i > 1 {
		instance, err := client.ExecSQl(query)
		if err == nil {
			err = instance.WaitForSuccess()
			if err == nil {
				res, err := instance.GetResult()
				if err == nil {
					for _, row := range res {
						if row.Content() != "" {
							return row.Content(), nil
						}
					}
					return "", errors.New("not able to parse row")
				}
			}
		}
		if err != nil {
			if strings.Contains(err.Error(), connErr) || strings.Contains(err.Error(), deadlineErr) {
				time.Sleep(200 * time.Millisecond)
				err = nil
				i--
				continue
			} else {
				return "", err
			}
		}
	}
	return "", errors.New("timeout")
}
