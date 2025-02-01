package tables

import (
	"github.com/spf13/cobra"
)

// NewTableCommand initializes command for project
func NewTableCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "table",
		Short:   "Commands that will let the user to operate on tables",
		Example: "opms mc table [sub-command]",
	}
	cmd.AddCommand(
		NewListCommand(),
		NewExistsCommand(),
		NewDropCommand(),
		NewReadCommand(),
	)

	return cmd
}
