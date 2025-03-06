package names

import (
	"errors"
	"strings"
)

func MapNames(projMap map[string]string, names []string) ([]string, error) {
	if len(projMap) == 0 {
		return names, nil
	}

	tableNames := make([]string, len(names))
	for i, field := range names {
		newName, err := MapName(projMap, field)
		if err != nil {
			return nil, err
		}

		tableNames[i] = newName
	}
	return tableNames, nil
}

func MapName(projMap map[string]string, name string) (string, error) {
	split := strings.Split(name, ".")
	if len(split) != 3 {
		return "", errors.New("invalid name format")
	}

	projDataset := split[0] + "." + split[1]
	if pd, ok := projMap[projDataset]; ok {
		n1 := pd + "." + split[2]
		return n1, nil
	}

	projName := split[0]
	if proj, ok := projMap[split[0]]; ok {
		projName = proj
	}

	newName := projName + "." + split[1] + "." + split[2]
	return newName, nil
}
