package optimus

import (
	"github.com/spf13/cobra"

	"github.com/sbchaos/opms/cmd/optimus/generate"
	"github.com/sbchaos/opms/lib/config"
)

func NewOptimusCommand(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "opt",
		Short:   "Commands for managing optimus related data",
		Example: "opms opt [sub-command]",
	}

	cmd.AddCommand(
		generate.NewGenerateCommand(cfg),
	)
	return cmd
}
