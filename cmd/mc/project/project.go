package project

import (
	"github.com/spf13/cobra"
)

// NewProjectCommand initializes command for project
func NewProjectCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "project",
		Short:   "Commands that will let the user to operate on project",
		Example: "opms project [sub-command]",
	}
	cmd.AddCommand(
		NewListCommand(),
	)
	return cmd
}
