package function

import (
	"github.com/spf13/cobra"

	"github.com/sbchaos/opms/lib/config"
)

// NewUDFCommand initializes command for udf
func NewUDFCommand(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "udf",
		Short:   "Commands that will let the user to operate on udf",
		Example: "opms mc udf [sub-command]",
	}
	cmd.AddCommand(
		NewCreateCommand(cfg),
		NewGetCommand(cfg),
		NewDropCommand(cfg),
		NewUpdateCommand(cfg),
	)
	return cmd
}
