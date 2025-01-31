package mc

import (
	"github.com/spf13/cobra"

	"github.com/sbchaos/opms/cmd/mc/project"
	"github.com/sbchaos/opms/cmd/mc/tables"
)

func NewMaxcomputeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "mc",
		Short:   "Commands for managing maxcompute related data",
		Example: "opms mc [sub-command]",
	}

	cmd.AddCommand(
		project.NewProjectCommand(),
		tables.NewTableCommand(),
	)
	return cmd
}
