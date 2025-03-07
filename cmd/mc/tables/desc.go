package tables

import (
	"fmt"

	"github.com/aliyun/aliyun-odps-go-sdk/odps"
	"github.com/spf13/cobra"

	mcc "github.com/sbchaos/opms/external/mc"
	"github.com/sbchaos/opms/lib/config"
	"github.com/sbchaos/opms/lib/names"
)

type descCommand struct {
	cfg *config.Config

	name string
}

// NewDescCommand returns data from the table
func NewDescCommand(cfg *config.Config) *cobra.Command {
	ec := &descCommand{cfg: cfg}

	cmd := &cobra.Command{
		Use:     "desc",
		Short:   "Describe details of a table in maxcompute",
		Example: "opms mc tables desc",
		RunE:    ec.RunE,
	}

	cmd.Flags().StringVarP(&ec.name, "name", "n", "", "Table name")
	cmd.MarkFlagRequired("name")
	return cmd
}

func (r *descCommand) RunE(_ *cobra.Command, _ []string) error {
	client, err := mcc.NewClientFromConfig(r.cfg)
	if err != nil {
		return err
	}

	tab, err := names.FromTableName(r.name)
	if err != nil {
		return err
	}

	client.SetDefaultProjectName(tab.Schema.ProjectID)
	client.SetCurrentSchemaName(tab.Schema.SchemaID)

	tabl := client.Table(tab.TableID)
	err = tabl.Load()
	if err != nil {
		return fmt.Errorf("failed to load table: %w", err)
	}

	printTable(tab, tabl)
	return nil
}

func printTable(ps names.Table, t *odps.Table) {
	fmt.Printf("Name:\t%s\n", ps.String())
	fmt.Printf("Type:\t%s\n", t.Type())
	fmt.Printf("Comment:\t%s\n", t.Comment())
	fmt.Printf("Last Changed Time:\t%s\n", t.LastModifiedTime())
	fmt.Printf("Size:\t%d\n", t.Size())
	fmt.Printf("Record Num:\t%d\n", t.RecordNum())
	if len(t.PartitionColumns()) > 0 {
		fmt.Printf("Partition Columns:\n")
		for _, c1 := range t.PartitionColumns() {
			fmt.Printf(" %s\n", c1.Name)
		}
	}
}
