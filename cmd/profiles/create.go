package profiles

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/sbchaos/opms/lib/cmdutil"
	"github.com/sbchaos/opms/lib/config"
	"github.com/sbchaos/opms/lib/keyring"
)

type createCommand struct {
	cfg *config.Config

	dynamic bool
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

	cmd.Flags().BoolVarP(&c.dynamic, "dynamic", "d", false, "Create a dynamic profile")

	return cmd
}

func (c *createCommand) RunE(_ *cobra.Command, _ []string) error {
	reader := bufio.NewReader(os.Stdin)
	profile := config.Profile{}

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

	if c.dynamic {
		profile.Dynamic = true
		profile.Creds = make(map[string]string)
		err := StoreDynamicCreds(reader, &profile)
		if err != nil {
			return err
		}
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
	fmt.Printf("Store for %s (yes/no): ", sys)
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

func StoreDynamicCreds(reader *bufio.Reader, p *config.Profile) error {
	fmt.Printf("Creating dynamic creds\n")

	fmt.Printf("Please provide system name (gcp|mc) :")
	sys, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	sys = strings.ToLower(strings.TrimSpace(sys))

	fmt.Printf("Please provide suffix for json files: ")
	suffix, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	suffix = strings.TrimSpace(suffix)

	fmt.Printf("\nProvide path for credentials folder: ")
	dir := getInput(reader)

	files, err := cmdutil.ListFiles(dir)
	if err != nil {
		return err
	}

	fileSuffix := suffix + ".json"
	for _, file := range files {
		if strings.HasSuffix(file, fileSuffix) {
			fmt.Printf("Found file %s\n", file)
			err = storeCred(file, sys, p, reader)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func storeCred(filename string, sys string, profile *config.Profile, reader *bufio.Reader) error {
	content, err := cmdutil.ReadFile(filename, nil)
	if err != nil {
		return err
	}

	var mapping map[string]string
	err = json.Unmarshal(content, &mapping)
	if err != nil {
		return err
	}

	proj, ok := mapping["project_id"]
	if !ok {
		proj, ok = mapping["project_name"]
		if !ok {
			proj = filepath.Base(filename)
		}
	}

	keyringKey := proj + "_" + profile.Name
	err = keyring.Set(keyringKey, string(content))
	if err != nil {
		return err
	}

	name := proj + "_" + sys
	profile.SetCred(name, keyringKey)

	fmt.Printf("Please provide projects for %s (, separated):", proj)
	otherProjs, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	others := strings.Split(otherProjs, ",")
	for _, other := range others {
		p1 := strings.TrimSpace(other)
		n1 := p1 + "_" + sys
		profile.SetCred(n1, keyringKey)
	}

	return nil
}
