package airflow

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/sbchaos/opms/external/airflow"
	"github.com/sbchaos/opms/lib/cmdutil"
	"github.com/sbchaos/opms/lib/color"
	"github.com/sbchaos/opms/lib/config"
	"github.com/sbchaos/opms/lib/printers/tree"
	"github.com/sbchaos/opms/lib/util"
)

type stuckCommand struct {
	cfg *config.Config

	name string

	authFile string
	afl      *airflow.Airflow

	reqCache map[string]bool
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
	s.reqCache = make(map[string]bool)

	treePrinter := tree.NewTreeWithAutoDetect[string]()
	errs := s.JobRunStatus(ctx, s.name, treePrinter.Root())

	treePrinter.Render(os.Stdout)

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
		n1 := FailureNode(name, "ErrorInFetch")
		parent.AddNode(n1)
		return errs
	}

	if dagRun == nil || len(dagRun.DagRuns) == 0 {
		errs = append(errs, fmt.Errorf("job runs not found for %s", name))
		n1 := FailureNode(name, "ErrorInFetch")
		parent.AddNode(n1)
		return errs
	}

	r1 := dagRun.DagRuns[0]
	node := Node(name, r1.State, r1.EndDate, 0)
	parent.AddNode(node)
	s.reqCache[name] = true

	instances, err := s.afl.TaskInstances(ctx, r1.DagID, r1.DagRunID)
	if err != nil {
		errs = append(errs, err)
		return errs
	}

	for _, t1 := range instances.TaskInstances {
		if strings.HasPrefix(t1.TaskDisplayName, "wait_") {
			lastIdx := strings.LastIndex(t1.TaskDisplayName, "-")
			upstreamName := t1.TaskDisplayName[5:lastIdx]

			if strings.EqualFold(t1.State, "success") {
				node.AddNode(SuccessNode(upstreamName, t1.EndDate))
				continue
			}

			_, ok := s.reqCache[upstreamName]
			if ok {
				node.AddNode(Node(upstreamName, "Repeat", time.Time{}, 0))
				continue
			}

			childErrs := s.JobRunStatus(ctx, upstreamName, node)
			if childErrs != nil {
				errs = append(errs, childErrs...)
			}

			continue
		}
		n1 := Node(t1.TaskDisplayName, t1.State, t1.EndDate, 0)
		node.AddNode(n1)
	}

	return errs
}

func Node(name, status string, end time.Time, col int) *tree.Node[string] {
	display := "Pending"
	if status != "" {
		display = status
	}
	if !end.IsZero() {
		display += " " + util.ToISO(end)
	}
	n1 := tree.NewNode(name, display)
	if col > 0 {
		n1.Color = col
	} else {
		if strings.EqualFold(status, "success") {
			n1.Color = color.Green
		} else if strings.EqualFold(status, "failure") {
			n1.Color = color.Red
		} else if strings.EqualFold(status, "up_for_retry") {
			n1.Color = color.Yellow
		} else if strings.EqualFold(status, "Repeat") {
			n1.Color = color.DarkGray
		}
	}
	return n1
}

func SuccessNode(name string, end time.Time) *tree.Node[string] {
	status := "Success"
	return Node(name, status, end, color.Green)
}

func FailureNode(name string, status string) *tree.Node[string] {
	return Node(name, status, time.Time{}, color.Red)
}
