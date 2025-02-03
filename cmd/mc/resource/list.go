package resource

import (
	"fmt"
	"os"

	"github.com/aliyun/aliyun-odps-go-sdk/odps"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	mcc "github.com/sbchaos/opms/external/mc"
	"github.com/sbchaos/opms/lib/config"
	"github.com/sbchaos/opms/lib/table"
	"github.com/sbchaos/opms/lib/term"
)

type listCommand struct {
	cfg *config.Config

	project string
	schema  string

	namePrefix   string
	resourceType string
}

func NewListCommand(cfg *config.Config) *cobra.Command {
	list := &listCommand{cfg: cfg}

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List the resources",
		Example: "opms mc resource list",
		RunE:    list.RunE,
	}

	cmd.Flags().StringVarP(&list.project, "project", "p", "", "Project")
	cmd.Flags().StringVarP(&list.schema, "schema", "s", "", "Schema")

	cmd.Flags().StringVarP(&list.namePrefix, "prefix", "n", "", "Resource name prefix")
	//cmd.Flags().StringVarP(&list.resourceType, "type", "t", "", "Resource type to query")
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

	client.SetDefaultProjectName(proj)
	client.SetCurrentSchemaName(schema)

	res := odps.NewResources(client)
	var filters []odps.RFileFunc = nil

	if r.namePrefix != "" {
		fnc := odps.ResourceFilter.NamePrefix(r.namePrefix)
		filters = append(filters, odps.RFileFunc(fnc))
	}

	t := term.FromEnv(0, 0)
	size, _, err := t.Size()
	if err != nil {
		size = 120
	}

	printer := table.New(os.Stdout, t.IsTerminalOutput(), size)
	fn := setupPrinter(printer)

	res.List(fn, filters...)

	err = printer.Render()
	if err != nil {
		return errors.Wrap(err, "failed to print resources")
	}
	return nil
}

func setupPrinter(printer table.Printer) func(*odps.Resource, error) {
	printer.AddHeader([]string{"Resource Name", "Type"})

	return func(r *odps.Resource, err error) {
		if err != nil {
			fmt.Printf("[ERROR] %s\n", err)
			return
		}

		printer.AddField(r.Name())
		printer.AddField(odps.ResourceTypeToStr(r.ResourceType()))
		printer.EndRow()
	}
}
