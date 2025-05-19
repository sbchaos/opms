package airflow

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/sbchaos/opms/external/airflow"
	"github.com/sbchaos/opms/lib/cmdutil"
	"github.com/sbchaos/opms/lib/config"
	"github.com/sbchaos/opms/lib/printers/tree"
	"github.com/sbchaos/opms/lib/util"
)

type stuckCommand struct {
	cfg *config.Config

	name string

	authFile string
	afl      *airflow.Airflow

	level int
}

func NewStuckCommand(cfg *config.Config) *cobra.Command {
	stuck := &stuckCommand{cfg: cfg}

	cmd := &cobra.Command{
		Use:     "stuck",
		Short:   "Find the deep upstream reason for stuck job",
		Example: "opms airflow stuck",
		RunE:    stuck.RunE,
	}

	cmd.Flags().StringVarP(&stuck.name, "name", "n", "", "Name of job")
	cmd.Flags().StringVarP(&stuck.authFile, "auth-file", "a", "", "Authentication json path")

	return cmd
}

func (s *stuckCommand) RunE(_ *cobra.Command, _ []string) error {
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

	afl := airflow.NewAirflow(auth)
	s.afl = afl

	root := tree.NewNode("Root", "Waiting")
	errs := s.JobRunStatus(ctx, s.name, root)

	fmt.Println(root.String())

	if len(errs) > 0 {
		for _, e1 := range errs {
			fmt.Println(e1)
		}
	}
	return nil
}

func (s *stuckCommand) JobRunStatus(ctx context.Context, name string, parent *tree.Node[string]) []error {
	errs := make([]error, 0)
	if parent.Level() > 15 {
		return nil
	}

	query := airflow.JobRunsCriteria{
		Name:        name,
		OnlyLastRun: true,
	}

	dagRun, err := s.afl.FetchJobRunBatch(ctx, &query)
	if err != nil {
		errs = append(errs, err)
		parent.AddChild(name, "ErrorInFetch")
		return errs
	}

	if dagRun == nil || len(dagRun.DagRuns) == 0 {
		errs = append(errs, fmt.Errorf("job runs not found for %s", name))
		parent.AddChild(name, "EmptyJobRun")
		return errs
	}

	r1 := dagRun.DagRuns[0]
	state := r1.State
	if !r1.EndDate.IsZero() {
		state += " " + util.ToISO(r1.EndDate)
	}
	node := tree.NewNode(name, state)
	parent.AddNode(node)

	instances, err := s.afl.TaskInstances(ctx, r1.DagID, r1.DagRunID)
	if err != nil {
		errs = append(errs, err)
		return errs
	}

	for _, t1 := range instances.TaskInstances {
		if strings.HasPrefix(t1.TaskDisplayName, "wait_") {
			lastIdx := strings.LastIndex(t1.TaskDisplayName, "-")
			upstreamName := t1.TaskDisplayName[5:lastIdx]

			if !strings.EqualFold(t1.State, "success") {
				childErrs := s.JobRunStatus(ctx, upstreamName, node)
				if childErrs != nil {
					errs = append(errs, childErrs...)
				}
			} else {
				status := fmt.Sprintf("success %s", util.ToISO(t1.EndDate))
				node.AddChild(upstreamName, status)
			}

			continue
		}
		status := "Waiting"
		if t1.State != "" {
			status = t1.State
			if !t1.EndDate.IsZero() {
				status += " " + util.ToISO(t1.EndDate)
			}
		}

		node.AddChild(t1.TaskDisplayName, status)
	}

	return errs
}
