package spec

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/sbchaos/opms/cmd/optimus/internal/io"
	"github.com/sbchaos/opms/lib/config"
)

type Spec struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
}

func (s Spec) SpecName() string {
	return s.Name
}

type duplicatesCommand struct {
	cfg *config.Config

	dir string
}

func NewDuplicatesCommand(cfg *config.Config) *cobra.Command {
	duplicates := &duplicatesCommand{cfg: cfg}

	cmd := &cobra.Command{
		Use:     "duplicates",
		Short:   "find duplicates",
		Example: "opms opt spec duplicates",
		RunE:    duplicates.RunE,
	}

	cmd.Flags().StringVarP(&duplicates.dir, "folder-path", "d", ".", "dir path")
	return cmd
}

func (r *duplicatesCommand) RunE(_ *cobra.Command, _ []string) error {
	jobNameMapping := map[string][]string{}
	resourceNameMapping := map[string][]string{}

	err := io.Walk[Spec](r.dir, jobNameMapping, resourceNameMapping)
	if err != nil {
		fmt.Printf("Unable to walk dir %s: %s", r.dir, err)
	}

	for name, paths := range jobNameMapping {
		if len(paths) > 1 {
			fmt.Printf("Found duplicates Jobs for %s\n", name)
			for _, path := range paths {
				fmt.Printf("\t%s\n", path)
			}
		}
	}

	for name, paths := range resourceNameMapping {
		if len(paths) > 1 {
			fmt.Printf("Found duplicates resources for %s\n", name)
			for _, path := range paths {
				fmt.Printf("\t%s\n", path)
			}
		}
	}

	return nil
}
