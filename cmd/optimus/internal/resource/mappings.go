package resource

import (
	"fmt"
	"strings"

	"github.com/sbchaos/opms/lib/names"
)

const (
	DecimalType = "DECIMAL"
)

type MappedExtTable struct {
	Et      *ExternalTable
	Table   names.Table
	OldName string
}

func MapExternalTable(spec *ExternalTable, projectMapping, typeMapping map[string]string) (*MappedExtTable, error) {
	schema, err := MapSchema(spec.Schema, typeMapping)
	if err != nil {
		return nil, err
	}

	oldName := spec.Name
	tab, err := names.MapName(projectMapping, oldName)
	if err != nil {
		return nil, err
	}

	mcET := ExternalTable{
		Name:        tab.TableID,
		Database:    tab.Schema.SchemaID,
		Project:     tab.Schema.ProjectID,
		Description: spec.Description,
		Schema:      schema,
		Source:      MapExternalSourceConfig(spec.Source),
	}

	mapped := &MappedExtTable{
		Et:      &mcET,
		Table:   tab,
		OldName: oldName,
	}

	return mapped, nil
}

func MapExternalSourceConfig(source *ExternalSource) *ExternalSource {
	m1 := ExternalSource{}
	if source == nil {
		return &m1
	}

	m1.SourceURIs = source.SourceURIs
	m1.SourceType = strings.ToUpper(source.SourceType)
	serde := map[string]string{
		"odps.sql.text.schema.mismatch.mode": "truncate",
		"odps.sql.text.option.flush.header":  "true",
	}

	rng, ok := source.Config["range"]
	if ok {
		m1.Range = fmt.Sprintf("%v", rng)
	}

	rws, ok := source.Config["skip_leading_rows"]
	if ok {
		serde["odps.text.option.header.lines.count"] = fmt.Sprintf("%v", rws)
	}
	m1.SerdeProperties = serde
	return &m1
}

func MapSchema(schema Schema, typeMapping map[string]string) (Schema, error) {
	var mcSchema []*Field
	if len(schema) == 0 {
		return nil, nil
	}

	for _, f1 := range schema {
		if f1.Type == "" && f1.Name == "" {
			continue
		}

		t2, ok := typeMapping[strings.ToUpper(f1.Type)]
		if !ok {
			// Use same type in case of no mapping
			t2 = f1.Type
		}

		mF1 := Field{
			Name: f1.Name,
			Type: t2,
		}
		if strings.EqualFold(t2, DecimalType) {
			mF1.Decimal = &Decimal{
				Precision: 38,
				Scale:     9,
			}
		}
		if f1.Description != "" {
			mF1.Description = f1.Description
		}
		mcSchema = append(mcSchema, &mF1)
	}
	return mcSchema, nil
}
