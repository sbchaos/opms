package tables

import (
	"github.com/spf13/cobra"

	"github.com/sbchaos/opms/lib/config"
)

// NewTableCommand initializes command for project
func NewTableCommand(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "table",
		Short:   "Commands that will let the user to operate on tables",
		Example: "opms bq table [sub-command]",
	}
	cmd.AddCommand(
		NewCountCommand(cfg),
	)

	return cmd
}
