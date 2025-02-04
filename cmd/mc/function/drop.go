package function

import (
	"github.com/aliyun/aliyun-odps-go-sdk/odps"
	"github.com/spf13/cobra"

	mcc "github.com/sbchaos/opms/external/mc"
	"github.com/sbchaos/opms/lib/config"
)

type dropCommand struct {
	cfg *config.Config

	project string
	schema  string

	name string
}

func NewDropCommand(cfg *config.Config) *cobra.Command {
	list := &dropCommand{cfg: cfg}

	cmd := &cobra.Command{
		Use:     "drop",
		Short:   "Drop the function",
		Example: "opms mc udf drop",
		RunE:    list.RunE,
	}

	cmd.Flags().StringVarP(&list.project, "project", "p", "", "Project")
	cmd.Flags().StringVarP(&list.schema, "schema", "s", "", "Schema")

	cmd.Flags().StringVarP(&list.name, "name", "n", "", "Function name prefix")
	cmd.MarkFlagRequired("name")
	return cmd
}

func (r *dropCommand) RunE(_ *cobra.Command, _ []string) error {
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

	functions := odps.NewFunctions(client)
	err = functions.Delete(r.name)
	if err != nil {
		return err
	}

	return nil
}
