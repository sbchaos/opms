package resource

import (
	"errors"
	"fmt"

	"github.com/mitchellh/mapstructure"
)

type YamlSpec struct {
	Version int                    `yaml:"version"`
	Name    string                 `yaml:"name"`
	Type    string                 `yaml:"type"`
	Labels  map[string]string      `yaml:"labels"`
	Spec    map[string]interface{} `yaml:"spec"`
	Path    string                 `yaml:"-"`
}

func ConvertSpecTo[T any](ys *YamlSpec) (*T, error) {
	var spec T
	if err := mapstructure.Decode(ys.Spec, &spec); err != nil {
		msg := fmt.Sprintf("%s: not able to decode spec for %s", err, ys.Name)
		return nil, errors.New(msg)
	}
	return &spec, nil
}
