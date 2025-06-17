package io

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func ReadSpec[T any](filePath string) (T, error) {
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
