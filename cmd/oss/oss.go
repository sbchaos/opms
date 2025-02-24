package oss

import (
	"github.com/spf13/cobra"

	"github.com/sbchaos/opms/lib/config"
)

func NewOSSCommand(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "oss",
		Short:   "Commands for managing oss",
		Example: "opms oss [sub-command]",
	}

	cmd.AddCommand(
		NewReadCommand(cfg),
	)
	return cmd
}
