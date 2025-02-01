package mc

import (
	"errors"
	"strings"
)

type ProjectSchema struct {
	ProjectName string
	SchemaName  string
}

func (ps ProjectSchema) String() string {
	return ps.ProjectName + "." + ps.SchemaName
}

func (ps ProjectSchema) Table(name string) string {
	return ps.ProjectName + "." + ps.ProjectName + "." + name
}

func SplitParts(name string) (ProjectSchema, string, error) {
	parts := strings.Split(name, ".")
	if len(parts) < 3 {
		return ProjectSchema{}, "", errors.New("invalid table name")
	}

	ps := ProjectSchema{
		ProjectName: parts[0],
		SchemaName:  parts[1],
	}
	return ps, parts[2], nil
}
