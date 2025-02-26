package names

import (
	"strings"
)

func MapNames(projMap, datasetMap map[string]string, names []string) ([]string, error) {
	if len(projMap) == 0 && len(datasetMap) == 0 {
		return names, nil
	}

	var tableNames []string

	for _, field := range names {
		split := strings.Split(field, ".")

		projName := split[0]
		if proj, ok := projMap[split[0]]; ok {
			projName = proj
		}

		schemaName := split[1]
		if schema, ok := datasetMap[schemaName]; ok {
			schemaName = schema
		}

		newName := projName + "." + schemaName + "." + split[2]
		tableNames = append(tableNames, newName)
	}
	return tableNames, nil
}
