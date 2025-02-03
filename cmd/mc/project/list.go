package project

import (
	"fmt"

	"github.com/aliyun/aliyun-odps-go-sdk/odps"
	"github.com/spf13/cobra"

	mcc "github.com/sbchaos/opms/external/mc"
	"github.com/sbchaos/opms/lib/config"
)

type listCommand struct {
	cfg *config.Config

	namePrefix string
}

// NewListCommand initializes command to list the projects
// Does not work reliably when one account used across projects
func NewListCommand(cfg *config.Config) *cobra.Command {
	list := &listCommand{cfg: cfg}

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List the registered projects",
		Example: "opms mc project list",
		RunE:    list.RunE,
	}

	cmd.Flags().StringVarP(&list.namePrefix, "prefix", "n", "", "Project name prefix")
	return cmd
}

func (r *listCommand) RunE(_ *cobra.Command, _ []string) error {
	client, err := mcc.NewClientFromConfig(r.cfg)
	if err != nil {
		return err
	}

	pjs := client.Projects()
	var filters []odps.PFilterFunc = nil
	if r.namePrefix != "" {
		filters = append(filters, odps.ProjectFilter.NamePrefix(r.namePrefix))
	}

	projs, err := pjs.List(filters...)
	if err != nil {
		return fmt.Errorf("list projects error: %w", err)
	}

	printProjs(projs)
	return nil
}

func printProjs(projs []*odps.Project) {
	for i, proj := range projs {
		fmt.Printf("%d %s %s %s\n", i+1, proj.Name(), proj.Type(), proj.LastModifiedTime())
	}
}
