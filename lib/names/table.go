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

func (s Schema) String() string {
	return fmt.Sprintf("%s.%s", s.ProjectID, s.SchemaID)
}

func (s Schema) TableName(name string) string {
	return fmt.Sprintf("%s.%s.%s", s.ProjectID, s.SchemaID, name)
}

func (t Table) String() string {
	return fmt.Sprintf("%s.%s.%s", t.Schema.ProjectID, t.Schema.SchemaID, t.TableID)
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

func TableWithSchema(schema Schema, tableID string) Table {
	return Table{
		Schema:  schema,
		TableID: tableID,
	}
}
