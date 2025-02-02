package project

import (
	"github.com/spf13/cobra"

	"github.com/sbchaos/opms/lib/config"
)

// NewProjectCommand initializes command for project
func NewProjectCommand(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "project",
		Short:   "Commands that will let the user to operate on project",
		Example: "opms project [sub-command]",
	}
	cmd.AddCommand(
		NewListCommand(cfg),
	)
	return cmd
}
