package tables

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/aliyun/aliyun-odps-go-sdk/odps"
	"github.com/spf13/cobra"

	mcc "github.com/sbchaos/opms/external/mc"
	"github.com/sbchaos/opms/lib/cmdutil"
	"github.com/sbchaos/opms/lib/table"
	"github.com/sbchaos/opms/lib/term"
)

type ProjectSchema struct {
	ProjectName string
	SchemaName  string
}

type existsCommand struct {
	credStr string

	name     string
	fileName string
}

// NewExistsCommand checks if the tables exist
func NewExistsCommand() *cobra.Command {
	ec := &existsCommand{}

	cmd := &cobra.Command{
		Use:     "exists",
		Short:   "Show if the table exists in maxcompute",
		Example: "opms mc tables exists",
		RunE:    ec.RunE,
	}

	cmd.Flags().StringVarP(&ec.credStr, "creds", "c", "", "Credentials in json format")

	cmd.Flags().StringVarP(&ec.name, "name", "n", "", "Table name")
	cmd.Flags().StringVarP(&ec.fileName, "filename", "f", "", "Filename with list of tables, - for stdin")
	return cmd
}

func (r *existsCommand) RunE(_ *cobra.Command, _ []string) error {
	t := term.FromEnv(0, 0)
	size, _, err := t.Size()
	if err != nil {
		size = 120
	}

	client, err := mcc.NewClient(r.credStr)
	if err != nil {
		return err
	}

	mapping := make(map[ProjectSchema][]string)
	if r.name != "" {
		parts, err := splitName(r.name)
		if err != nil {
			return err
		}

		mapping[ProjectSchema{
			ProjectName: parts[0],
			SchemaName:  parts[1],
		}] = []string{parts[2]}

		if r.fileName != "" {
			return errors.New("--filename flag cannot be used along with name")
		}
	}

	if r.fileName != "" {
		content, err := cmdutil.ReadFile(r.fileName, os.Stdin)
		if err != nil {
			return err
		}
		fields := strings.Fields(string(content))
		for _, field := range fields {
			nameParts, err := splitName(field)
			if err != nil {
				fmt.Printf("ignoring invalid table name %s\n", field)
			}

			ps := ProjectSchema{
				ProjectName: nameParts[0],
				SchemaName:  nameParts[1],
			}

			mapping[ps] = append(mapping[ps], nameParts[2])
		}
	}

	printer := table.New(os.Stdout, t.IsTerminalOutput(), size)

	for ps, tables := range mapping {
		client.SetDefaultProjectName(ps.ProjectName)
		client.SetCurrentSchemaName(ps.SchemaName)

		lenTables := len(tables)
		if lenTables < 10 {
			forEvery(printer, client.Tables(), ps, tables)
		} else {
			for i := 0; i < lenTables; i = i + 50 {
				end := min(i+50, lenTables)
				forN(50, printer, client.Tables(), ps, tables[i:end])
			}
		}
	}

	err = printer.Render()
	if err != nil {
		return fmt.Errorf("failed to print table: %w", err)
	}
	return nil
}

func forEvery(printer table.Printer, tabs *odps.Tables, ps ProjectSchema, tables []string) {
	for _, t1 := range tables {
		tres, err := tabs.BatchLoadTables([]string{t1})
		if err != nil {
			failureStatus(printer, ps, t1)
		} else {
			successStatus(printer, ps, tres)
		}
	}
}

func forN(step int, printer table.Printer, tabs *odps.Tables, ps ProjectSchema, tables []string) {
	lenTables := len(tables)
	for i := 0; i < lenTables; i = i + step {
		end := min(i+step, lenTables)
		loadTables, err := tabs.BatchLoadTables(tables[i:end])
		if err == nil {
			// If no error mark all as success
			successStatus(printer, ps, loadTables)
			continue
		} else {
			fmt.Printf("error in batch: %s, trying smaller size", err)
			forEvery(printer, tabs, ps, tables)
		}
	}
}

func splitName(name string) ([]string, error) {
	parts := strings.Split(name, ".")
	if len(parts) < 3 {
		return parts, errors.New("invalid table name")
	}

	return parts, nil
}

func failureStatus(printer table.Printer, ps ProjectSchema, table string) {
	printer.AddField(" ❌ ")
	printer.AddField(ps.ProjectName + "." + ps.SchemaName + "." + table)
	printer.EndRow()
}

func successStatus(printer table.Printer, ps ProjectSchema, tables []*odps.Table) {
	prefix := ps.ProjectName + "." + ps.SchemaName + "."
	printer.AddHeader([]string{"Exists", "Table Name"})

	for _, t1 := range tables {
		printer.AddField(" ✅ ")
		printer.AddField(prefix + t1.Name())
		printer.EndRow()
	}
}

// p_gopay_id_mart.gopay_consolidated.external_payment_dashboard_service_type_reference_flag
