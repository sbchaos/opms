package airflow

import (
	"github.com/spf13/cobra"

	"github.com/sbchaos/opms/lib/config"
)

func NewAirflowCommand(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "airflow",
		Aliases: []string{"afl"},
		Short:   "Commands for managing airflow related data",
		Example: "opms afl [sub-command]",
	}

	cmd.AddCommand(
		NewStatusCommand(cfg),
		NewWatchCommand(cfg),
		NewRunsCommand(cfg),
	)
	return cmd
}
