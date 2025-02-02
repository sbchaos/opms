package profiles

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/sbchaos/opms/lib/config"
	"github.com/sbchaos/opms/lib/jsonpretty"
)

// NewShowProfileCommand returns data from the table
func NewShowProfileCommand(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "show",
		Short:   "Show a profile",
		Example: "opms profile show <profile>",
		RunE:    ShowProfile,
	}

	return cmd
}

func ShowProfile(_ *cobra.Command, _ []string) error {
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
				bytes, _ := json.Marshal(p)
				jsonpretty.Format(os.Stdout, strings.NewReader(string(bytes)), " ", true)

				return nil
			}
		}
		fmt.Printf("Please choose a profile: ")
	}
}
