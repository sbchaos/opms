package spec

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/sbchaos/opms/cmd/optimus/internal/job"
	"github.com/sbchaos/opms/lib/config"
)

type endDateCommand struct {
	cfg *config.Config

	dir     string
	endDate string

	// Filters
	jobNames      string
	taskNames     string
	deleteEndDate string
}

func NewEndDateCommand(cfg *config.Config) *cobra.Command {
	endDate := &endDateCommand{cfg: cfg}

	cmd := &cobra.Command{
		Use:     "endDate",
		Short:   "update the endDate for job spec",
		Example: "opms opt spec endDate",
		RunE:    endDate.RunE,
	}

	cmd.Flags().StringVarP(&endDate.dir, "folder-path", "f", ".", "dir path")
	cmd.Flags().StringVarP(&endDate.jobNames, "job-names", "j", "", "Jobs to target")
	cmd.Flags().StringVarP(&endDate.taskNames, "task-names", "t", "", "Comma separated task names")
	cmd.Flags().StringVarP(&endDate.endDate, "end-date", "e", "", "End date for job spec")
	cmd.Flags().StringVarP(&endDate.deleteEndDate, "remove-end-date", "d", "", "Remove the end date")
	return cmd
}

func (r *endDateCommand) RunE(_ *cobra.Command, _ []string) error {
	tasks := map[string]bool{}
	checkTasks := false
	if r.taskNames != "" {
		for _, t := range strings.Split(r.taskNames, ",") {
			t1 := strings.TrimSpace(t)
			tasks[t1] = true
		}
		checkTasks = len(tasks) > 0
	}

	jobs := map[string]bool{}
	checkJobs := false
	if r.jobNames != "" {
		for _, n := range strings.Split(r.jobNames, ",") {
			n1 := strings.TrimSpace(n)
			jobs[n1] = true
		}
		checkJobs = len(jobs) > 0
	}

	walker := func(path string, d fs.DirEntry, _ error) error {
		if d.IsDir() {
			return nil
		}

		fileName := filepath.Base(path)
		if fileName != "job.yaml" && fileName != "job.yml" {
			return nil
		}

		spec, err := readSpec[job.YamlSpec](path)
		if err != nil {
			fmt.Printf("Unable to read spec for %s: %s", path, err)
		}

		if checkJobs {
			if _, ok := jobs[spec.Name]; !ok {
				return nil
			}
		}

		if checkTasks {
			_, ok := tasks[spec.Task.Name]
			if !ok {
				return nil
			}
		}

		if spec.Schedule.EndDate != "" {
			fmt.Printf("End date for job %s: %s\n", spec.Name, spec.Schedule.EndDate)
		} else {
			fmt.Printf("Processing %s\n", spec.Name)
			spec.Schedule.EndDate = r.endDate
			writeSpec(path, spec)
		}

		return nil
	}

	err := filepath.WalkDir(r.dir, walker)
	if err != nil {
		fmt.Printf("Unable to walk dir %s: %s", r.dir, err)
	}

	return nil
}

func writeSpec(filePath string, spec job.YamlSpec) error {
	fileSpec, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating spec under [%s]: %w", filePath, err)
	}
	indent := 2
	encoder := yaml.NewEncoder(fileSpec)
	encoder.SetIndent(indent)
	return encoder.Encode(spec)
}
