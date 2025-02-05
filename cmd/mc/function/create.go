package function

import (
	"github.com/aliyun/aliyun-odps-go-sdk/odps"
	"github.com/spf13/cobra"

	mcc "github.com/sbchaos/opms/external/mc"
	"github.com/sbchaos/opms/lib/config"
)

type createCommand struct {
	cfg *config.Config

	project string
	schema  string

	name      string
	classPath string
	resource  string
}

func NewCreateCommand(cfg *config.Config) *cobra.Command {
	list := &createCommand{cfg: cfg}

	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Create the function",
		Example: "opms mc udf create",
		RunE:    list.RunE,
	}

	cmd.Flags().StringVarP(&list.project, "project", "p", "", "Project")
	cmd.Flags().StringVarP(&list.schema, "schema", "s", "", "Schema")

	cmd.Flags().StringVarP(&list.name, "name", "n", "", "Function name prefix")
	cmd.Flags().StringVarP(&list.classPath, "classpath", "c", "", "Resource type to query")
	cmd.Flags().StringVarP(&list.resource, "res", "r", "", "Resource type to query")
	cmd.MarkFlagRequired("name")
	cmd.MarkFlagRequired("classpath")
	cmd.MarkFlagRequired("res")
	return cmd
}

func (r *createCommand) RunE(_ *cobra.Command, _ []string) error {
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

	fun := odps.
		NewFunctionBuilder().
		Name(r.name).
		SchemaName(schema).
		ClassPath(r.classPath).
		Resources([]string{r.resource}).
		Build()

	functions := odps.NewFunctions(client)
	err = functions.Create(proj, schema, fun)
	if err != nil {
		return err
	}

	return nil
}
