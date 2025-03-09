package vars

import (
	"github.com/spf13/cobra"

	"github.com/sbchaos/opms/lib/config"
)

func NewVarsCommand(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "var",
		Short:   "Commands for managing variables in profile",
		Example: "opms profile var [sub-command]",
	}

	cmd.AddCommand(
		NewAddCommand(cfg),
		NewDeleteCommand(cfg),
	)
	return cmd
}
