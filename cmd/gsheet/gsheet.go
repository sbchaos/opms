package gsheet

import (
	"github.com/spf13/cobra"

	"github.com/sbchaos/opms/lib/config"
)

func NewGsheetsCommand(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "gsheets",
		Short:   "Commands for managing gsheets",
		Example: "opms gsheets [sub-command]",
	}

	cmd.AddCommand(
		NewReadCommand(cfg),
	)
	return cmd
}
