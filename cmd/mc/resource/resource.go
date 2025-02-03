package resource

import (
	"github.com/spf13/cobra"

	"github.com/sbchaos/opms/lib/config"
)

// NewResourceCommand initializes command for resource
func NewResourceCommand(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "resource",
		Short:   "Commands that will let the user to operate on resource",
		Example: "opms resource [sub-command]",
	}
	cmd.AddCommand(
		NewListCommand(cfg),
	)
	return cmd
}
