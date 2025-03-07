package names

func MapNames(projMap map[string]string, names []string) ([]Table, error) {
	tables := make([]Table, len(names))
	if len(projMap) == 0 {
		for i, name := range names {
			t1, err := FromTableName(name)
			if err != nil {
				return nil, err
			}
			tables[i] = t1
		}
	}

	for i, field := range names {
		newName, err := MapName(projMap, field)
		if err != nil {
			return nil, err
		}

		tables[i] = newName
	}
	return tables, nil
}

func MapName(projMap map[string]string, name string) (Table, error) {
	t1, err := FromTableName(name)
	if err != nil {
		return Table{}, err
	}

	schm := t1.Schema.String()
	if pd, ok := projMap[schm]; ok {
		schema, err := FromSchemaName(pd)
		if err != nil {
			return Table{}, err
		}
		return TableWithSchema(schema, t1.TableID), nil
	}

	projName := t1.Schema.ProjectID
	if proj, ok := projMap[projName]; ok {
		projName = proj
	}

	tab := NewTable(projName, t1.Schema.SchemaID, t1.TableID)
	return tab, nil
}
