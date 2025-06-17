package plan

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/sbchaos/opms/cmd/optimus/internal/conf"
	"github.com/sbchaos/opms/cmd/optimus/internal/io"
	"github.com/sbchaos/opms/cmd/optimus/internal/plan"
	"github.com/sbchaos/opms/lib/cmdutil"
	"github.com/sbchaos/opms/lib/config"
)

const outputFileName = "resource.json"

type Spec struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
}

func (s Spec) SpecName() string {
	return s.Name
}

type planCommand struct {
	cfg *config.Config

	jobsFile       string
	jobNamesFilter []string

	resourcesFile       string
	resourceNamesFilter []string

	namespace string
}

func NewPlanCommand(cfg *config.Config) *cobra.Command {
	p1 := &planCommand{cfg: cfg}

	cmd := &cobra.Command{
		Use:     "plan",
		Short:   "gen plan",
		Example: "opms opt gen plan",
		RunE:    p1.RunE,
	}

	cmd.Flags().StringVarP(&p1.namespace, "namespace", "n", "", "Namespace to generate plan for")
	cmd.Flags().StringVarP(&p1.jobsFile, "jobs", "J", "", "Filename with list of jobs, - for stdin")
	cmd.Flags().StringVarP(&p1.resourcesFile, "resources", "R", "", "Filename with list of resources, - for stdin")
	return cmd
}

func (r *planCommand) RunE(_ *cobra.Command, _ []string) error {
	spec, err := io.ReadSpec[conf.ClientConfig](conf.DefaultFilename)
	if err != nil {
		return fmt.Errorf("read client config: %w", err)
	}

	if r.jobsFile != "" {
		fields, err := cmdutil.ReadLines(r.jobsFile, os.Stdin)
		if err != nil {
			return err
		}
		r.jobNamesFilter = fields
	}

	if r.resourcesFile != "" {
		fields, err := cmdutil.ReadLines(r.resourcesFile, os.Stdin)
		if err != nil {
			return err
		}
		r.resourceNamesFilter = fields
	}

	p1 := &plan.Plan{
		ProjectName: spec.Project.Name,
		Job: plan.OpsPerNS[plan.JobPlan]{
			Create: make(map[string][]plan.JobPlan),
		},
		Resource: plan.OpsPerNS[plan.ResourcePlan]{
			Create: make(map[string][]plan.ResourcePlan),
		},
	}

	if r.namespace != "" {
		for _, namespace := range spec.Namespaces {
			if namespace.Name == r.namespace {
				r.RunForNamespace(namespace, p1)
			}
		}
	} else {
		for _, ns := range spec.Namespaces {
			r.RunForNamespace(ns, p1)
		}
	}

	bytes, err := json.MarshalIndent(p1, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal plan error: %w", err)
	}

	fmt.Println("Writing optimus plan")
	return cmdutil.WriteFile(outputFileName, bytes)
}

func (r *planCommand) RunForNamespace(ns *conf.Namespace, toPlan *plan.Plan) {
	jobNameMapping := map[string][]string{}
	resourceNameMapping := map[string][]string{}
	// read all job specs
	err := io.Walk[Spec](ns.Job.Path, jobNameMapping, resourceNameMapping)
	if err != nil {
		fmt.Printf("Unable to walk dir %s: %s", ns.Job.Path, err)
	}

	if len(ns.Datastore) > 0 {
		err2 := io.Walk[Spec](ns.Datastore[0].Path, jobNameMapping, resourceNameMapping)
		if err2 != nil {
			fmt.Printf("Unable to walk dir %s: %s", ns.Datastore[0].Path, err)
		}
	}

	if len(r.jobNamesFilter) > 0 {
		filteredJobs := map[string][]string{}
		for _, field := range r.jobNamesFilter {
			paths, ok := jobNameMapping[field]
			if ok {
				filteredJobs[field] = paths
			}
		}
		jobNameMapping = filteredJobs
	}

	if len(r.resourceNamesFilter) > 0 {
		filteredResources := map[string][]string{}
		for _, field := range r.resourceNamesFilter {
			paths, ok := jobNameMapping[field]
			if ok {
				filteredResources[field] = paths
			}
		}
		resourceNameMapping = filteredResources
	}

	var jobPlans []plan.JobPlan
	for jobName, paths := range jobNameMapping {
		if len(paths) > 1 {
			fmt.Printf("Job[%s] has multiple files %s\n", jobName, paths)
		}

		j1 := plan.JobPlan{
			Name:         jobName,
			OldNamespace: nil,
			Path:         paths[0],
		}
		jobPlans = append(jobPlans, j1)
	}

	var resourcePlans []plan.ResourcePlan
	for resName, paths := range resourceNameMapping {
		if len(paths) > 1 {
			fmt.Printf("Resource[%s] has multiple files %s\n", resName, paths)
		}

		r1 := plan.ResourcePlan{
			Name:         resName,
			OldNamespace: nil,
			Path:         paths[0],
			Datastore:    ns.Datastore[0].Type,
		}
		resourcePlans = append(resourcePlans, r1)
	}

	if len(jobPlans) > 0 {
		toPlan.Job.Create[ns.Name] = jobPlans
	}
	if len(resourcePlans) > 0 {
		toPlan.Resource.Create[ns.Name] = resourcePlans
	}
}
