package profiles

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/sbchaos/opms/lib/config"
)

type edit struct {
	cfg *config.Config
}

// NewEditProfileCommand returns data from the table
func NewEditProfileCommand(cfg *config.Config) *cobra.Command {
	s := &edit{cfg: cfg}
	cmd := &cobra.Command{
		Use:     "edit",
		Short:   "Edit a profile",
		Example: "opms profile edit <profile>",
		RunE:    s.RunE,
	}

	return cmd
}

func (s *edit) RunE(_ *cobra.Command, _ []string) error {
	fmt.Printf("Available Profile:\n")
	for i, p := range s.cfg.AvailableProfiles {
		fmt.Printf("%d. %s\n", i+1, p.Name)
	}

	updatedProfile := s.chooseProfile()

	var newProfiles []config.Profile
	for _, pr := range s.cfg.AvailableProfiles {
		if pr.Name != updatedProfile.Name {
			newProfiles = append(newProfiles, pr)
		}
	}
	newProfiles = append(newProfiles, *updatedProfile)
	s.cfg.AvailableProfiles = newProfiles
	s.cfg.CurrentProfile = updatedProfile.Name

	err := config.Write(s.cfg)
	if err != nil {
		fmt.Printf("Error writing config: %s\n", err)
	}

	return nil
}

func (s *edit) chooseProfile() *config.Profile {
	reader := bufio.NewReader(os.Stdin)
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
			editProfile(&p, reader)
			return &p
		}

		for _, p := range s.cfg.AvailableProfiles {
			if p.Name == nameInput {
				editProfile(&p, reader)
				return &p
			}
		}
		fmt.Printf("Please choose a profile: ")
	}
}

func editProfile(p *config.Profile, reader *bufio.Reader) {
	if p.GCPCred == "" {
		key, err := StoreCredsFor(reader, p.Name, "GCP")
		if err == nil && key != "" {
			p.GCPCred = key
			fmt.Printf("Stored creds for GCP\n")
		}
	}

	if p.MCCred == "" {
		key, err := StoreCredsFor(reader, p.Name, "Maxcompute")
		if err == nil && key != "" {
			p.MCCred = key
			fmt.Printf("Stored creds for Maxcompute\n")
		}
	}
}
