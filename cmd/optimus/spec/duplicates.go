package spec

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/sbchaos/opms/lib/config"
)

type Spec struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
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

	walker := func(path string, d fs.DirEntry, _ error) error {
		if d.IsDir() {
			return nil
		}

		jobYaml := false
		resourceYaml := false

		fileName := filepath.Base(path)
		switch fileName {
		case "job.yaml", "job.yml":
			jobYaml = true
		case "resource.yaml", "resource.yml":
			resourceYaml = true
		}

		if !jobYaml && !resourceYaml {
			return nil
		}

		spec, err := readSpec[Spec](path)
		if err != nil {
			fmt.Printf("Unable to read spec for %s: %s", path, err)
		}

		if spec.Name == "" {
			return nil
		}

		if jobYaml {
			jobNameMapping[spec.Name] = append(jobNameMapping[spec.Name], path)
		} else if resourceYaml {
			resourceNameMapping[spec.Name] = append(resourceNameMapping[spec.Name], path)
		}

		return nil
	}

	err := filepath.WalkDir(r.dir, walker)
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

func readSpec[T any](filePath string) (T, error) {
	var spec T
	fileSpec, err := os.Open(filePath)
	if err != nil {
		return spec, fmt.Errorf("error opening spec under [%s]: %w", filePath, err)
	}
	defer fileSpec.Close()

	if err = yaml.NewDecoder(fileSpec).Decode(&spec); err != nil {
		return spec, fmt.Errorf("error decoding spec under [%s]: %w", filePath, err)
	}

	return spec, nil
}
