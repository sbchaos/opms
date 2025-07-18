package generate

import (
	"github.com/spf13/cobra"

	"github.com/sbchaos/opms/cmd/optimus/generate/external_table"
	"github.com/sbchaos/opms/cmd/optimus/generate/plan"
	"github.com/sbchaos/opms/lib/config"
)

func NewGenerateCommand(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "generate",
		Aliases: []string{"gen"},
		Short:   "Commands for generating optimus specs",
		Example: "opms optimus generate [sub-command]",
	}

	cmd.AddCommand(
		external_table.NewExternalTableCommand(cfg),
		plan.NewPlanCommand(cfg),
	)
	return cmd
}
