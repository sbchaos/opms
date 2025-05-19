package tables

import (
	"fmt"
	"os"

	"github.com/aliyun/aliyun-odps-go-sdk/odps"
	"github.com/spf13/cobra"

	mcc "github.com/sbchaos/opms/external/mc"
	"github.com/sbchaos/opms/lib/config"
	"github.com/sbchaos/opms/lib/printers/table"
	"github.com/sbchaos/opms/lib/term"
)

type listCommand struct {
	cfg *config.Config

	project string
	schema  string

	namePrefix string
	tableType  string
}

// NewListCommand initializes command to list the projects
// Does not work reliably when one account used across projects
func NewListCommand(cfg *config.Config) *cobra.Command {
	list := &listCommand{cfg: cfg}

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List the tables",
		Example: "opms mc tables list",
		RunE:    list.RunE,
	}

	cmd.Flags().StringVarP(&list.project, "project", "p", "", "Project")
	cmd.Flags().StringVarP(&list.schema, "schema", "s", "", "Schema")

	cmd.Flags().StringVarP(&list.namePrefix, "prefix", "n", "", "Table name prefix")
	cmd.Flags().StringVarP(&list.tableType, "type", "t", "", "Table type to query, eg MANAGED_TABLE, VIRTUAL_TABLE, EXTERNAL_TABLE")
	return cmd
}

func (r *listCommand) RunE(_ *cobra.Command, _ []string) error {
	client, err := mcc.NewClientFromConfig(r.cfg)
	if err != nil {
		return err
	}

	schema := "default"
	if r.schema != "" {
		schema = r.schema
	}

	proj := client.DefaultProjectName()
	if r.project != "" {
		proj = r.project
	}

	tables := odps.NewTables(client, proj, schema)
	var filters []odps.TFilterFunc = nil

	if r.namePrefix != "" {
		filters = append(filters, odps.TableFilter.NamePrefix(r.namePrefix))
	}

	if r.tableType != "" {
		typ := odps.TableTypeFromStr(r.tableType)
		if typ == odps.TableTypeUnknown {
			return fmt.Errorf("invalid table type: %s", r.tableType)
		}
		filters = append(filters, odps.TableFilter.Type(typ))
	}

	t := term.FromEnv(0, 0)
	size, _ := t.Size(120)

	printer := table.New(os.Stdout, t.IsTerminalOutput(), size)
	fn := setupPrinter(printer)

	tables.List(fn, filters...)

	err = printer.Render()
	if err != nil {
		return fmt.Errorf("failed to print table: %w", err)
	}
	return nil
}

func setupPrinter(printer table.Printer) func(*odps.Table, error) {
	printer.AddHeader([]string{"Table Name", "Type", "Last Update"})

	return func(table *odps.Table, err error) {
		if err != nil {
			fmt.Printf("[ERROR] %s\n", err)
			return
		}

		printer.AddField(table.Name())
		printer.AddField(table.Type().String())
		printer.AddField(table.LastModifiedTime().String())
		printer.EndRow()
	}
}
