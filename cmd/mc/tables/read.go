package tables

import (
	"fmt"

	"github.com/aliyun/aliyun-odps-go-sdk/odps"
	"github.com/spf13/cobra"

	mcc "github.com/sbchaos/opms/external/mc"
)

type readCommand struct {
	credStr string

	name  string
	limit int
}

// NewReadCommand returns data from the table
func NewReadCommand() *cobra.Command {
	ec := &readCommand{}

	cmd := &cobra.Command{
		Use:     "read",
		Short:   "Read contents of a table in maxcompute",
		Example: "opms mc tables read",
		RunE:    ec.RunE,
	}

	cmd.Flags().StringVarP(&ec.credStr, "creds", "c", "", "Credentials in json format")

	cmd.Flags().StringVarP(&ec.name, "name", "n", "", "Table name")
	cmd.Flags().IntVarP(&ec.limit, "limit", "l", 1, "Limit number of results")
	cmd.MarkFlagRequired("name")
	return cmd
}

func (r *readCommand) RunE(_ *cobra.Command, _ []string) error {
	client, err := mcc.NewClient(r.credStr)
	if err != nil {
		return err
	}

	ps, name, err := mcc.SplitParts(r.name)
	if err != nil {
		return err
	}

	client.SetDefaultProjectName(ps.ProjectName)
	client.SetCurrentSchemaName(ps.SchemaName)

	tab := client.Table(name)
	err = tab.Load()
	if err != nil {
		return fmt.Errorf("failed to load table: %w", err)
	}

	printTable(ps, tab)

	//query := "select * from " + ps.Table(tab.Name())
	//if limit := r.limit; limit > 0 {
	//	query += " limit " + strconv.Itoa(limit)
	//}
	//query += ";"
	//
	//instance, err := tab.ExecSql("query", query)
	//if err != nil {
	//	return err
	//}
	//
	//err = instance.WaitForSuccess()
	//if err != nil {
	//	return err
	//}
	//
	//fmt.Println(instance.TaskResults())
	return nil
}

func printTable(ps mcc.ProjectSchema, t *odps.Table) {
	fmt.Printf("Name:\t%s\n", ps.Table(t.Name()))
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
