package io

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/sbchaos/opms/cmd/optimus/internal/job"
)

func WriteSpec(filePath string, spec job.YamlSpec) error {
	fileSpec, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating spec under [%s]: %w", filePath, err)
	}
	indent := 2
	encoder := yaml.NewEncoder(fileSpec)
	encoder.SetIndent(indent)
	return encoder.Encode(spec)
}
