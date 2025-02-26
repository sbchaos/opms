package names

import "fmt"

func GroupTableNames(names []string) map[Schema][]string {
	groups := make(map[Schema][]string)

	for _, name := range names {
		tn, err := FromTableName(name)
		if err != nil {
			fmt.Printf("Ignoring %s: with err: %s\n", name, err)
			continue
		}

		groups[tn.Schema] = append(groups[tn.Schema], tn.TableID)
	}
	return groups
}
