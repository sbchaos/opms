package airflow

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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

type watchCommand struct {
	cfg *config.Config

	fileName string

	authFile string
	auth     airflow.Auth

	interval int
	client   *airflow.Client
	mu       *sync.Mutex
}

func NewWatchCommand(cfg *config.Config) *cobra.Command {
	watch := &watchCommand{cfg: cfg}

	cmd := &cobra.Command{
		Use:     "watch",
		Short:   "Watch the status jobs on airflow",
		Example: "opms airflow watch",
		RunE:    watch.RunE,
	}

	cmd.Flags().StringVarP(&watch.fileName, "filename", "f", "", "Filename with list of jobs to enable/disable")
	cmd.Flags().StringVarP(&watch.authFile, "auth-file", "a", "", "Authentication json path")
	cmd.Flags().IntVarP(&watch.interval, "interval", "i", 5, "Refresh interval in seconds")

	return cmd
}

func (s *watchCommand) RunE(_ *cobra.Command, _ []string) error {
	var auth airflow.Auth
	if s.authFile == "" {
		return fmt.Errorf("--auth-file is required")
	}

	err := cmdutil.ReadJsonFile(s.authFile, os.Stdin, &auth)
	if err != nil {
		return err
	}
	s.auth = auth

	if s.fileName == "" {
		return fmt.Errorf("--filename is required")
	}

	jobNames, err := cmdutil.ReadLines(s.fileName, os.Stdin)
	if err != nil {
		return err
	}

	client := airflow.NewAirflowClient()
	s.client = client
	s.mu = &sync.Mutex{}

	ctx, cancelFunc := context.WithTimeout(context.Background(), timeout)
	defer cancelFunc()

	t := term.FromEnv(0, 0)
	size, _ := t.Size(120)
	printer := table.New(os.Stdout, t.IsTerminalOutput(), size)
	printer.AddHeader([]string{"Job", "Status", "Interval", "Next_Run"})

	tasks := make([]func() pool.JobResult[string], len(jobNames))
	for i, t1 := range jobNames {
		tasks[i] = func() pool.JobResult[string] {
			err := s.getJobStatus(ctx, t1, printer)
			return pool.JobResult[string]{
				Output: t1,
				Err:    err,
			}
		}
	}

	for {
		printer.Clear()
		fmt.Println("\033[H\033c")
		outchan := pool.RunWithWorkers(5, tasks)
		printer.Render()

		for out := range outchan {
			if out.Err != nil {
				fmt.Printf("Error for job [%s]:%s\n", out.Output, out.Err)
			}
		}

		time.Sleep(time.Duration(s.interval) * time.Second)
	}
	return nil
}

func (s *watchCommand) getJobStatus(ctx context.Context, jobName string, printer table.Printer) error {
	status := "Failed"

	req := airflow.Request{
		Path:   fmt.Sprintf(dagURL, jobName),
		Method: http.MethodGet,
	}
	var resp airflow.DAGObj
	content, err := s.client.Invoke(ctx, req, s.auth)
	if err == nil {
		err = json.Unmarshal(content, &resp)
		if err == nil {
			if resp.IsPaused {
				status = "Paused"
			} else {
				status = "UnPaused"
			}
		} else {
			err = fmt.Errorf("invalid response from Airflow: %v", err)
		}
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	printer.AddField(jobName)
	printer.AddField(status)
	printer.AddField(resp.ScheduleInterval.Value)
	printer.AddField(resp.NextDagRun)
	printer.EndRow()
	return err
}
