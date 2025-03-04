package sql

import (
	"github.com/spf13/cobra"

	"github.com/sbchaos/opms/lib/config"
)

func NewSQLCommand(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "sql",
		Short:   "Commands that will provide sql related functions",
		Example: "opms mc sql [sub-command]",
	}
	cmd.AddCommand(
		NewRunSQLCommand(cfg),
	)

	return cmd
}
