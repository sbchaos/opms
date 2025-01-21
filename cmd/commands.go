package cmd

import (
	cli "github.com/spf13/cobra"
)

// New constructs the 'root' command. It houses all other sub commands
// default output of logging should go to stdout
// interactive output like progress bars should go to stderr
// unless the stdout/err is a tty, colors/progressbar should be disabled
func New() *cli.Command {
	cmd := &cli.Command{
		Use:          "opms <command> <subcommand> [flags]",
		Long:         "",
		SilenceUsage: true,
		Example: `
				$ opms project list
			`,
		Annotations: map[string]string{
			"group:core": "true",
			"help:learn": `
				Use 'opms <command> <subcommand> --help' for more information about a command.
			`,
		},
	}

	cmd.AddCommand()

	return cmd
}
