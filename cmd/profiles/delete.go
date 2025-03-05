package profiles

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/sbchaos/opms/lib/config"
	"github.com/sbchaos/opms/lib/keyring"
)

type delete struct {
	cfg *config.Config
}

// NewDeleteProfileCommand returns data from the table
func NewDeleteProfileCommand(cfg *config.Config) *cobra.Command {
	s := &delete{cfg: cfg}
	cmd := &cobra.Command{
		Use:     "delete",
		Short:   "Delete a profile",
		Example: "opms profile delete <profile>",
		RunE:    s.RunE,
	}

	return cmd
}

func (s delete) RunE(_ *cobra.Command, _ []string) error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Available Profile:\n")
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
		found := false
		for _, p := range s.cfg.AvailableProfiles {
			if p.Name == nameInput {
				found = true
			}
		}

		if !found {
			fmt.Printf("Please choose correct profile: ")
			continue
		}

		var restOfProfile []config.Profile
		for _, p := range s.cfg.AvailableProfiles {
			if p.Name == nameInput {
				if p.GCPCred != "" {
					keyring.Delete(p.GCPCred)
				}
				if p.MCCred != "" {
					keyring.Delete(p.GCPCred)
				}
				if p.Airflow != "" {
					keyring.Delete(p.GCPCred)
				}

				for _, v := range p.Creds {
					keyring.Delete(v)
				}
			} else {
				restOfProfile = append(restOfProfile, p)
			}
		}

		s.cfg.AvailableProfiles = restOfProfile
		config.Write(s.cfg)

		fmt.Printf("Profile %s deleted\n", nameInput)
		return nil
	}
}
