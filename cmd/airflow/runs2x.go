package airflow

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/sbchaos/opms/external/airflow"
	"github.com/sbchaos/opms/lib/cmdutil"
	"github.com/sbchaos/opms/lib/config"
)

type runsCommand struct {
	cfg *config.Config

	name string

	authFile string

	onlyLastRun bool
	startTime   string
	endTime     string
	status      string
}

func NewRunsCommand(cfg *config.Config) *cobra.Command {
	runs := &runsCommand{cfg: cfg}

	cmd := &cobra.Command{
		Use:     "runs",
		Short:   "Update the runs of jobs on airflow",
		Example: "opms airflow runs",
		RunE:    runs.RunE,
	}

	cmd.Flags().StringVarP(&runs.name, "name", "n", "", "Name of job")
	cmd.Flags().StringVarP(&runs.authFile, "auth-file", "a", "", "Authentication json path")
	cmd.Flags().BoolVarP(&runs.onlyLastRun, "last", "l", false, "Get only last run")
	cmd.Flags().StringVarP(&runs.startTime, "start", "s", "", "Start time for interval")
	cmd.Flags().StringVarP(&runs.endTime, "end", "e", "", "End time for interval")
	cmd.Flags().StringVarP(&runs.status, "status", "t", "", "Status of job")

	return cmd
}

func (s *runsCommand) RunE(_ *cobra.Command, _ []string) error {
	var auth airflow.Auth
	if s.authFile == "" {
		return fmt.Errorf("--auth-file is required")
	}

	err := cmdutil.ReadJsonFile(s.authFile, os.Stdin, &auth)
	if err != nil {
		return err
	}

	if s.name == "" {
		return fmt.Errorf("--name is required")
	}

	ctx, cancelFunc := context.WithTimeout(context.Background(), timeout)
	defer cancelFunc()

	query := airflow.JobRunsCriteria{
		Name: s.name,
	}
	if s.onlyLastRun {
		query.OnlyLastRun = true
	} else {
		if s.startTime != "" {
			start, err := time.Parse(time.RFC3339, s.startTime)
			if err != nil {
				return err
			}
			query.StartDate = start
		}
		if s.endTime != "" {
			end, err := time.Parse(time.RFC3339, s.endTime)
			if err != nil {
				return err
			}
			query.EndDate = end
		}

		if s.status == "" {
			query.Filter = []string{s.status}
		}
	}

	afl := airflow.NewAirflow(auth)
	runs, err := afl.FetchJobRunBatch(ctx, &query)
	if err != nil {
		return err
	}

	fmt.Printf("%s: Runs[%d]\n", s.name, len(runs.DagRuns))
	for _, r1 := range runs.DagRuns {
		fmt.Printf("\nStatus: %s\n", r1.State)
		fmt.Printf("Logical/Execution Date: %s\n", r1.LogicalDate.Format(time.RFC3339))
		fmt.Printf("Execution: Start[%s] -> END[%s]\n", r1.StartDate.Format(time.RFC3339), r1.EndDate.Format(time.RFC3339))
		fmt.Printf("Interval:  Start[%s] -> END[%s]\n", r1.DataIntervalStart.Format(time.RFC3339), r1.DataIntervalEnd.Format(time.RFC3339))
	}

	return nil
}
