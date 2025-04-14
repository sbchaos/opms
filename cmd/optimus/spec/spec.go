package spec

import (
	"github.com/spf13/cobra"

	"github.com/sbchaos/opms/lib/config"
)

func NewSpecCommand(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "spec",
		Short:   "Commands for managing optimus related data",
		Example: "opms spec [sub-command]",
	}

	cmd.AddCommand(
		NewDuplicatesCommand(cfg),
		NewEndDateCommand(cfg),
	)
	return cmd
}
