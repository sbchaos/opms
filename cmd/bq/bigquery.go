package bq

import (
	"github.com/spf13/cobra"

	"github.com/sbchaos/opms/cmd/bq/tables"
	"github.com/sbchaos/opms/lib/config"
)

func NewBigQueryCommand(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "bq",
		Short:   "Commands for managing bigquery related data",
		Example: "opms bq [sub-command]",
	}

	cmd.AddCommand(
		tables.NewTableCommand(cfg),
	)
	return cmd
}
