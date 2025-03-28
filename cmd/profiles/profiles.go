package profiles

import (
	"github.com/spf13/cobra"

	"github.com/sbchaos/opms/cmd/profiles/vars"
	"github.com/sbchaos/opms/lib/config"
)

func NewProfilesCommand(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "profile",
		Short:   "Commands for managing profiles",
		Example: "opms profile [sub-command]",
	}

	cmd.AddCommand(
		NewCreateProfileCommand(cfg),
		NewShowProfileCommand(cfg),
		NewUseProfileCommand(cfg),
		NewDeleteProfileCommand(cfg),
		NewEditProfileCommand(cfg),
		vars.NewVarsCommand(cfg),
	)
	return cmd
}
