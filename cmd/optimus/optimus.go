package optimus

import (
	"github.com/spf13/cobra"

	"github.com/sbchaos/opms/cmd/optimus/generate"
	_map "github.com/sbchaos/opms/cmd/optimus/map"
	"github.com/sbchaos/opms/cmd/optimus/spec"
	"github.com/sbchaos/opms/lib/config"
)

func NewOptimusCommand(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "opt",
		Short:   "Commands for managing optimus related data",
		Example: "opms opt [sub-command]",
	}

	cmd.AddCommand(
		generate.NewGenerateCommand(cfg),
		_map.NewMapNameCommand(cfg),
		spec.NewSpecCommand(cfg),
	)
	return cmd
}
