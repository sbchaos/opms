package vars

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/sbchaos/opms/lib/config"
)

type deleteCommand struct {
	cfg *config.Config

	scope string
	key   string
}

func NewDeleteCommand(cfg *config.Config) *cobra.Command {
	c := &deleteCommand{cfg: cfg}
	cmd := &cobra.Command{
		Use:     "delete",
		Short:   "Delete variable from profile",
		Example: "opms profile var delete",
		RunE:    c.RunE,
	}

	cmd.Flags().StringVarP(&c.scope, "scope", "s", "", "Add variable to a scope")
	cmd.Flags().StringVarP(&c.key, "key", "k", "", "Key for the var")

	return cmd
}

func (c *deleteCommand) RunE(_ *cobra.Command, _ []string) error {
	p1 := c.cfg.GetCurrentProfile()

	if c.key == "" {
		return errors.New("you must provide a key to delete")
	}
	key := c.key

	if c.scope != "" {
		key = c.scope + ":" + key
	}

	delete(p1.Variables, key)
	err := config.Write(c.cfg)
	if err != nil {
		fmt.Printf("Error writing config: %s\n", err)
	}

	return nil
}
