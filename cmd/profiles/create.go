package profiles

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/sbchaos/opms/lib/cmdutil"
	"github.com/sbchaos/opms/lib/config"
	"github.com/sbchaos/opms/lib/keyring"
)

type createCommand struct {
	cfg *config.Config
}

// NewCreateProfileCommand returns data from the table
func NewCreateProfileCommand(cfg *config.Config) *cobra.Command {
	c := &createCommand{cfg: cfg}
	cmd := &cobra.Command{
		Use:     "new",
		Short:   "Create a new profile",
		Example: "opms profile new",
		RunE:    c.RunE,
	}

	return cmd
}

func (c createCommand) RunE(_ *cobra.Command, _ []string) error {
	reader := bufio.NewReader(os.Stdin)
	profile := config.Profiles{}

	fmt.Printf("Please provide profile name: ")
	nameInput, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	nameInput = strings.TrimSpace(nameInput)
	profile.Name = nameInput

	for _, pf := range c.cfg.AvailableProfiles {
		if pf.Name == nameInput {
			return errors.New("profile already exists")
		}
	}

	key, err := StoreCredsFor(reader, nameInput, "GCP")
	if err == nil && key != "" {
		profile.GCPCred = key
		fmt.Printf("Stored creds for GCP\n")
	}

	key, err = StoreCredsFor(reader, nameInput, "Maxcompute")
	if err == nil && key != "" {
		profile.MCCred = key
		fmt.Printf("Stored creds for Maxcompute\n")
	}

	c.cfg.AvailableProfiles = append(c.cfg.AvailableProfiles, profile)
	c.cfg.CurrentProfile = profile.Name

	err = config.Write(c.cfg)
	if err != nil {
		fmt.Printf("Error writing config: %s\n", err)
	}

	return nil
}

func getYesNo(reader *bufio.Reader) bool {
	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("not able to read value, please try again")
			continue
		}

		input = strings.TrimSpace(input)
		if input == "yes" || input == "y" {
			return true
		}
		if input == "no" || input == "n" {
			return false
		}
		fmt.Println("Please enter 'yes' or 'no'")
	}
}

func getInput(reader *bufio.Reader) string {
	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("not able to read value, please try again")
			continue
		}

		input = strings.TrimSpace(input)
		if input == "" {
			fmt.Println("Please enter a value, please try again")
			continue
		}
		return input
	}
}

func StoreCredsFor(reader *bufio.Reader, name, sys string) (string, error) {
	fmt.Printf("Want to continue(yes/no) for %s: ", sys)
	proceed := getYesNo(reader)
	if !proceed {
		return "", nil
	}

	fmt.Printf("\nProvide path for credentials file: ")
	jpath := getInput(reader)
	bytes, err := cmdutil.ReadFile(jpath, os.Stdin)
	if err != nil {
		fmt.Printf("Error reading json file: %s\n", err)
	}

	key := name + "_" + strings.ToLower(sys)
	err = keyring.Set(key, string(bytes))
	if err != nil {
		fmt.Printf("Error setting keyring: %s\n", err)
		return "", err
	}

	return key, nil
}
