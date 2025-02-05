package verify

import (
	"github.com/spf13/cobra"

	"github.com/sbchaos/opms/lib/config"
)

// NewVerifyCommand initializes command for project
func NewVerifyCommand(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "verify",
		Short:   "Commands that will let the user to operate on verify",
		Example: "opms mc verify [sub-command]",
	}
	cmd.AddCommand(
		NewExternalTableCommand(cfg),
	)

	return cmd
}
