package profiles

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/sbchaos/opms/lib/config"
	"github.com/sbchaos/opms/lib/jsonpretty"
	"github.com/sbchaos/opms/lib/keyring"
)

type show struct {
	cfg *config.Config
}

var showKeyring = false

// NewShowProfileCommand returns data from the table
func NewShowProfileCommand(cfg *config.Config) *cobra.Command {
	s := &show{cfg: cfg}
	cmd := &cobra.Command{
		Use:     "show",
		Short:   "Show a profile",
		Example: "opms profile show <profile>",
		RunE:    s.RunE,
	}

	cmd.Flags().BoolVarP(&showKeyring, "keyring", "k", false, "-k show the value in keyring")
	cmd.Flags().MarkHidden("keyring")
	return cmd
}

func (s show) RunE(_ *cobra.Command, _ []string) error {
	fmt.Printf("\nKeyring %v\n", showKeyring)
	fmt.Println("")

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

		num, err := strconv.Atoi(nameInput)
		if err == nil {
			p := s.cfg.AvailableProfiles[num-1]
			showProfile(p)
			if showKeyring {
				showFromKeyring(p)
			}
			return nil
		}

		for _, p := range s.cfg.AvailableProfiles {
			if p.Name == nameInput {
				showProfile(p)
				if showKeyring {
					showFromKeyring(p)
				}
				return nil
			}
		}
		fmt.Printf("Please choose a profile: ")
	}
}

func showProfile(p config.Profiles) {
	bytes, _ := json.Marshal(p)
	jsonpretty.Format(os.Stdout, strings.NewReader(string(bytes)), " ", true)
}

func showFromKeyring(p config.Profiles) {
	if p.GCPCred != "" {
		fmt.Printf("GCP:\n")
		val, err := keyring.Get(p.GCPCred)
		if err != nil {
			fmt.Printf("error: %s", err)
		}
		jsonpretty.Format(os.Stdout, strings.NewReader(val), " ", true)
	}
	if p.MCCred != "" {
		fmt.Printf("Maxcompute:\n")
		val, err := keyring.Get(p.MCCred)
		if err != nil {
			fmt.Printf("error: %s", err)
		}
		jsonpretty.Format(os.Stdout, strings.NewReader(val), " ", true)
	}
}
