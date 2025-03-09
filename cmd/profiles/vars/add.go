package vars

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/sbchaos/opms/lib/config"
)

type addCommand struct {
	cfg *config.Config

	scope string
	key   string
	value string
}

func NewAddCommand(cfg *config.Config) *cobra.Command {
	c := &addCommand{cfg: cfg}
	cmd := &cobra.Command{
		Use:     "add",
		Short:   "Add a variable to a profile",
		Example: "opms profile var add",
		RunE:    c.RunE,
	}

	cmd.Flags().StringVarP(&c.scope, "scope", "s", "", "Add variable to a scope")
	cmd.Flags().StringVarP(&c.key, "key", "k", "", "Key for the var")
	cmd.Flags().StringVarP(&c.value, "value", "v", "", "Value for the var")

	return cmd
}

func (c *addCommand) RunE(_ *cobra.Command, _ []string) error {
	p1 := c.cfg.GetCurrentProfile()

	if c.key == "" {
		return errors.New("you must provide a key to add to")
	}
	key := c.key

	if c.value == "" {
		return errors.New("you must provide a value")
	}

	if c.scope != "" {
		key = c.scope + ":" + key
	}

	if p1.Variables == nil {
		p1.Variables = make(map[string]interface{})
	}

	p1.SetVariable(key, c.value)
	err := config.Write(c.cfg)
	if err != nil {
		fmt.Printf("Error writing config: %s\n", err)
	}

	return nil
}
