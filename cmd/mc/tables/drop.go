package tables

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"

	mcc "github.com/sbchaos/opms/external/mc"
	"github.com/sbchaos/opms/lib/cmdutil"
	"github.com/sbchaos/opms/lib/config"
	"github.com/sbchaos/opms/lib/names"
	"github.com/sbchaos/opms/lib/table"
	"github.com/sbchaos/opms/lib/term"
)

type dropCommand struct {
	cfg *config.Config

	name     string
	fileName string
}

// NewDropCommand checks if the tables exist
func NewDropCommand(cfg *config.Config) *cobra.Command {
	ec := &dropCommand{cfg: cfg}

	cmd := &cobra.Command{
		Use:     "drop",
		Short:   "Drop if a table exists in maxcompute",
		Example: "opms mc tables drop",
		RunE:    ec.RunE,
	}

	cmd.Flags().StringVarP(&ec.name, "name", "n", "", "Table name")
	cmd.Flags().StringVarP(&ec.fileName, "filename", "f", "", "Filename with list of tables, - for stdin")
	return cmd
}

func (r *dropCommand) RunE(_ *cobra.Command, _ []string) error {
	t := term.FromEnv(0, 0)
	size, _ := t.Size(120)

	client, err := mcc.NewClientFromConfig(r.cfg)
	if err != nil {
		return err
	}

	var tableNames []string
	mapping := make(map[names.Schema][]string)
	if r.name != "" {
		tableNames = append(tableNames, r.name)

		if r.fileName != "" {
			return errors.New("--filename flag cannot be used along with name")
		}
	}

	if r.fileName != "" {
		lines, err := cmdutil.ReadLines(r.fileName, os.Stdin)
		if err != nil {
			return err
		}
		tableNames = lines
	}
	mapping = names.GroupTableNames(tableNames)

	printer := table.New(os.Stdout, t.IsTerminalOutput(), size)

	for ps, tables := range mapping {
		client.SetDefaultProjectName(ps.ProjectID)
		client.SetCurrentSchemaName(ps.SchemaID)

		for _, t1 := range tables {
			for {
				err = client.Tables().Delete(t1, true)
				if err != nil {
					if strings.Contains(err.Error(), connErr) {
						time.Sleep(500 * time.Millisecond)
						continue
					}

					dropFailure(printer, ps, t1)
				} else {
					dropSuccess(printer, ps, t1)
				}
				break
			}
		}
	}

	err = printer.Render()
	if err != nil {
		return fmt.Errorf("failed to print table: %w", err)
	}
	return nil
}

func dropFailure(printer table.Printer, n1 names.Schema, n2 string) {
	printer.AddField(" ❌ ")
	printer.AddField(n1.TableName(n2))
	printer.EndRow()
}

func dropSuccess(printer table.Printer, n1 names.Schema, n2 string) {
	printer.AddField(" ✅ ")
	printer.AddField(n1.TableName(n2))
	printer.EndRow()
}
