package _map

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/sbchaos/opms/lib/cmdutil"
	"github.com/sbchaos/opms/lib/config"
	"github.com/sbchaos/opms/lib/names"
)

type mapNameCommand struct {
	cfg *config.Config

	name      string
	namesFile string
	projMap   string
}

func NewMapNameCommand(cfg *config.Config) *cobra.Command {
	mapName := &mapNameCommand{cfg: cfg}

	cmd := &cobra.Command{
		Use:   "map_name",
		Short: "Map the table names using project",
		Example: `opms opt map_name
It can provide an mapping file in json format.
eg
{
"proj_a": "proj_b",
"proj_b.dataset_b": "proj_c.dataset_d"
}
`,
		RunE: mapName.RunE,
	}

	cmd.Flags().StringVarP(&mapName.name, "name", "n", "", "Table name")
	cmd.Flags().StringVarP(&mapName.namesFile, "names-file", "f", "", "Table names file")
	cmd.Flags().StringVarP(&mapName.projMap, "mapping", "m", "", "Mapping for project dataset")
	return cmd
}

func (r *mapNameCommand) RunE(_ *cobra.Command, _ []string) error {
	projectMapping := map[string]string{}
	if r.projMap != "" {
		err := cmdutil.ReadJsonFile(r.projMap, os.Stdin, &projectMapping)
		if err != nil {
			return err
		}
	}

	var tableNames []string
	if r.name != "" {
		tableNames = []string{r.name}
	}

	if r.namesFile != "" {
		lines, err := cmdutil.ReadLines(r.namesFile, os.Stdin)
		if err != nil {
			return err
		}
		tableNames = lines
	}

	mappedNames, err := names.MapNames(projectMapping, tableNames)
	if err != nil {
		return err
	}

	for _, n := range mappedNames {
		fmt.Println(n)
	}
	return nil
}
