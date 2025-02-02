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

type show struct {
	cfg *config.Config
}

// NewShowProfileCommand returns data from the table
func NewShowProfileCommand(cfg *config.Config) *cobra.Command {
	s := &show{cfg: cfg}
	cmd := &cobra.Command{
		Use:     "show",
		Short:   "Show a profile",
		Example: "opms profile show <profile>",
		RunE:    s.RunE,
	}

	return cmd
}

func (s show) RunE(_ *cobra.Command, _ []string) error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Available Profiles:\n")
	for i, p := range s.cfg.AvailableProfiles {
		fmt.Printf("%d. %s\n", i+1, p.Name)
	}

	fmt.Printf("Please choose a profile: ")
	for {
		nameInput, err := reader.ReadString('\n')
		if err != nil {
			continue
		}
		nameInput = strings.TrimSpace(nameInput)
		for _, p := range s.cfg.AvailableProfiles {
			if p.Name == nameInput {
				bytes, _ := json.Marshal(p)
				jsonpretty.Format(os.Stdout, strings.NewReader(string(bytes)), " ", true)

				return nil
			}
		}
		fmt.Printf("Please choose a profile: ")
	}
}
