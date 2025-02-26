package names

import (
	"fmt"
	"strings"
)

type Schema struct {
	ProjectID string
	SchemaID  string
}

type Table struct {
	Schema  Schema
	TableID string
}

func NewSchema(projectID, schemaID string) Schema {
	return Schema{
		ProjectID: projectID,
		SchemaID:  schemaID,
	}
}

func FromSchemaName(name string) (Schema, error) {
	parts := strings.SplitN(name, ".", 2)
	if len(parts) != 2 {
		return Schema{}, fmt.Errorf("invalid schema name: %s", name)
	}
	return NewSchema(parts[0], parts[1]), nil
}

func FromTableName(name string) (Table, error) {
	parts := strings.Split(name, ".")
	if len(parts) < 3 {
		return Table{}, fmt.Errorf("invalid table name %q", name)
	}

	return NewTable(parts[0], parts[1], parts[2]), nil
}

func NewTable(projectID, schemaID, tableID string) Table {
	return Table{
		Schema:  NewSchema(projectID, schemaID),
		TableID: tableID,
	}
}
