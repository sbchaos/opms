package drive

import (
	"github.com/spf13/cobra"

	"github.com/sbchaos/opms/lib/config"
)

func NewDriveCommand(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "drive",
		Short:   "Commands for managing drive data",
		Example: "opms drive [sub-command]",
	}

	cmd.AddCommand(
		NewDownloadCommand(cfg),
		NewSyncCommand(cfg),
		NewDeleteCommand(cfg),
	)
	return cmd
}
