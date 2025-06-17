package io

import (
	"fmt"
	"io/fs"
	"path/filepath"
)

type Named interface {
	SpecName() string
}

func Walk[T Named](dir string, jobMap map[string][]string, resourceMap map[string][]string) error {
	walker := func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Printf("Skipping %s, err: %s\n", path, err)
			return nil
		}

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

		spec, err := ReadSpec[T](path)
		if err != nil {
			fmt.Printf("Unable to read spec for %s: %s", path, err)
		}

		name := spec.SpecName()
		if name == "" {
			return nil
		}

		if jobYaml {
			jobMap[name] = append(jobMap[name], path)
		} else if resourceYaml {
			resourceMap[name] = append(resourceMap[name], path)
		}

		return nil
	}

	err := filepath.WalkDir(dir, walker)
	return err
}
