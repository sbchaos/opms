package profiles

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/sbchaos/opms/lib/config"
)

// NewUseProfileCommand returns data from the table
func NewUseProfileCommand(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "use",
		Short:   "Switch the profile currently in use",
		Example: "opms profile use <profile>",
		RunE:    UseProfile,
	}

	return cmd
}

func UseProfile(_ *cobra.Command, _ []string) error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Available Profiles:\n ")
	read, err := config.Read(config.DefaultConfig())
	if err != nil {
		fmt.Printf("Error reading config: %s\n", err)
	}

	for i, p := range read.AvailableProfiles {
		fmt.Printf("%d. %s\n", i+1, p.Name)
	}

	fmt.Printf("Please choose a profile: ")
	for {
		nameInput, err := reader.ReadString('\n')
		if err != nil {
			continue
		}
		nameInput = strings.TrimSpace(nameInput)
		for _, p := range read.AvailableProfiles {
			if p.Name == nameInput {
				read.CurrentProfile = p.Name
				config.Write(read)
				return nil
			}
		}
		fmt.Printf("Please choose a profile: ")
	}
}
