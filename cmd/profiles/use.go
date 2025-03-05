package profiles

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/sbchaos/opms/lib/config"
)

type useCmd struct {
	cfg *config.Config
}

// NewUseProfileCommand returns data from the table
func NewUseProfileCommand(cfg *config.Config) *cobra.Command {
	u := &useCmd{cfg: cfg}
	cmd := &cobra.Command{
		Use:     "use",
		Short:   "Switch the profile currently in use",
		Example: "opms profile use <profile>",
		RunE:    u.RunE,
	}

	return cmd
}

func (u useCmd) RunE(_ *cobra.Command, _ []string) error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Available Profile:\n")
	for i, p := range u.cfg.AvailableProfiles {
		fmt.Printf("%d. %s\n", i+1, p.Name)
	}

	fmt.Printf("Please choose a profile: ")
	for {
		nameInput, err := reader.ReadString('\n')
		if err != nil {
			continue
		}
		nameInput = strings.TrimSpace(nameInput)
		for _, p := range u.cfg.AvailableProfiles {
			if p.Name == nameInput {
				u.cfg.CurrentProfile = p.Name
				config.Write(u.cfg)
				return nil
			}
		}
		fmt.Printf("Please choose a profile: ")
	}
}
