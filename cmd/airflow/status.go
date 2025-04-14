package airflow

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cobra"

	"github.com/sbchaos/opms/external/airflow"
	"github.com/sbchaos/opms/lib/cmdutil"
	"github.com/sbchaos/opms/lib/config"
	"github.com/sbchaos/opms/lib/pool"
	"github.com/sbchaos/opms/lib/table"
	"github.com/sbchaos/opms/lib/term"
)

var timeout = time.Minute * 5
var dagURL = "api/v1/dags/%s"

type statusCommand struct {
	cfg *config.Config

	name     string
	fileName string

	status string

	authFile string
	auth     airflow.Auth

	workers int
	client  *airflow.Client
	mu      *sync.Mutex
}

func NewStatusCommand(cfg *config.Config) *cobra.Command {
	status := &statusCommand{cfg: cfg}

	cmd := &cobra.Command{
		Use:     "status",
		Short:   "Update the status of jobs on airflow",
		Example: "opms airflow status",
		RunE:    status.RunE,
	}

	cmd.Flags().StringVarP(&status.status, "status", "s", "disable", "enable/disable")
	cmd.Flags().StringVarP(&status.name, "name", "n", "", "Name of job to enable/disable")
	cmd.Flags().StringVarP(&status.fileName, "filename", "f", "", "Filename with list of jobs to enable/disable")
	cmd.Flags().StringVarP(&status.authFile, "auth-file", "a", "", "Authentication json path")
	cmd.Flags().IntVarP(&status.workers, "workers", "w", 1, "Number of parallel workers")

	return cmd
}

func (s *statusCommand) RunE(_ *cobra.Command, _ []string) error {
	var auth airflow.Auth
	if s.authFile == "" {
		return fmt.Errorf("--auth-file is required")
	}

	err := cmdutil.ReadJsonFile(s.authFile, os.Stdin, &auth)
	if err != nil {
		return err
	}
	s.auth = auth

	if s.name == "" && s.fileName == "" {
		return fmt.Errorf("--name or --filename is required")
	}

	jobNames := make([]string, 0)
	if s.name != "" {
		jobNames = append(jobNames, s.name)
	} else {
		lines, err := cmdutil.ReadLines(s.fileName, os.Stdin)
		if err != nil {
			return err
		}
		jobNames = lines
	}

	var data []byte
	switch strings.ToLower(s.status) {
	case "enabled":
		data = []byte(`{"is_paused": false}`)
	case "disabled":
		data = []byte(`{"is_paused": true}`)
	default:
		return errors.New("unknown status" + s.status)
	}

	ctx, cancelFunc := context.WithTimeout(context.Background(), timeout)
	defer cancelFunc()

	client := airflow.NewAirflowClient()
	s.client = client
	s.mu = &sync.Mutex{}

	t := term.FromEnv(0, 0)
	size, _ := t.Size(120)
	printer := table.New(os.Stdout, t.IsTerminalOutput(), size)
	printer.AddHeader([]string{"Job", "Status"})

	tasks := make([]func() pool.JobResult[string], len(jobNames))
	for i, t1 := range jobNames {
		tasks[i] = func() pool.JobResult[string] {
			err := s.updateJobState(ctx, t1, data, printer)
			return pool.JobResult[string]{
				Output: t1,
				Err:    err,
			}
		}
	}

	outchan := pool.RunWithWorkers(s.workers, tasks)
	printer.Render()

	for out := range outchan {
		if out.Err != nil {
			fmt.Printf("Error for job [%s]:%s\n", out.Output, out.Err)
		}
	}
	return nil
}

func (s *statusCommand) updateJobState(ctx context.Context, jobName string, data []byte, printer table.Printer) error {
	req := airflow.Request{
		Path:   fmt.Sprintf(dagURL, jobName),
		Method: http.MethodPatch,
		Body:   data,
	}
	_, err := s.client.Invoke(ctx, req, s.auth)

	s.mu.Lock()
	defer s.mu.Unlock()
	printer.AddField(jobName)
	if err == nil {
		printer.AddField("Success")
	} else {
		printer.AddField("Failed")
	}

	printer.EndRow()
	return err
}
