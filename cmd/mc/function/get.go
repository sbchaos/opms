package function

import (
	"fmt"

	"github.com/aliyun/aliyun-odps-go-sdk/odps"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	mcc "github.com/sbchaos/opms/external/mc"
	"github.com/sbchaos/opms/lib/config"
)

type getCommand struct {
	cfg *config.Config

	project string
	schema  string

	name string
}

func NewGetCommand(cfg *config.Config) *cobra.Command {
	list := &getCommand{cfg: cfg}

	cmd := &cobra.Command{
		Use:     "get",
		Short:   "Get the function",
		Example: "opms mc udf get",
		RunE:    list.RunE,
	}

	cmd.Flags().StringVarP(&list.project, "project", "p", "", "Project")
	cmd.Flags().StringVarP(&list.schema, "schema", "s", "", "Schema")

	cmd.Flags().StringVarP(&list.name, "name", "n", "", "Function name prefix")
	//cmd.Flags().StringVarP(&list.resourceType, "type", "t", "", "Resource type to query")
	return cmd
}

func (r *getCommand) RunE(_ *cobra.Command, _ []string) error {
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

	if r.name == "" {
		return errors.New("must specify a function name")
	}

	functions := odps.NewFunctions(client)
	function, err := functions.Get(r.name)
	if err != nil {
		return err
	}

	printFunction(function)

	return nil
}

func printFunction(function *odps.Function) {
	fmt.Println("Function Details:")
	fmt.Printf("  Name: %s\n", function.Name())
	fmt.Printf("  IsSQL: %t\n", function.IsSQLFunction())
	fmt.Printf("  Language: %s\n", function.ProgramLanguage())
	fmt.Println("  Resources:")
	for i, s := range function.Resources() {
		fmt.Printf("    %d. %s\n", i, s)
	}
}
