package macro

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/sbchaos/opms/lib/cmdutil"
	"github.com/sbchaos/opms/lib/compiler"
	"github.com/sbchaos/opms/lib/config"
)

var (
	start       = time.Date(2024, 10, 01, 0, 0, 0, 0, time.UTC)
	end         = time.Date(2024, 10, 02, 0, 0, 0, 0, time.UTC)
	executedAt  = time.Date(2023, 10, 03, 0, 0, 0, 0, time.UTC)
	scheduledAt = time.Date(2024, 10, 02, 0, 0, 0, 0, time.UTC)
)

type macroCommand struct {
	cfg *config.Config

	compiler *compiler.Engine

	fileName string
	dirName  string

	output string
}

func NewMacroCommand(cfg *config.Config) *cobra.Command {
	macro := &macroCommand{cfg: cfg}

	cmd := &cobra.Command{
		Use:     "macro",
		Short:   "Replace the macros in a query",
		Example: "opms opt macro",
		RunE:    macro.RunE,
		PreRunE: macro.PreRunE,
	}

	cmd.Flags().StringVarP(&macro.fileName, "query", "q", "", "Query file name, - for stdin")
	cmd.Flags().StringVarP(&macro.dirName, "query-dir", "d", "", "Directory with queries")
	cmd.Flags().StringVarP(&macro.output, "output-dir", "o", "", "Directory for output")

	return cmd
}

func (m *macroCommand) PreRunE(_ *cobra.Command, _ []string) error {
	var err error
	m.compiler = compiler.NewEngine()

	return err
}

func (m *macroCommand) RunE(_ *cobra.Command, _ []string) error {
	names := make([]string, 0)
	if m.fileName != "" {
		names = append(names, m.fileName)
	}

	if m.dirName != "" && len(names) == 0 {
		n1, err := cmdutil.ListFiles(m.dirName)
		if err != nil {
			return err
		}

		names = n1
	}

	timeConfigs := getTimeConfigs(start, end, executedAt, scheduledAt)

	if err := os.MkdirAll(m.output, os.ModePerm); err != nil {
		return err
	}
	for i, name := range names {
		processed, err := m.processWithTemplate(name, timeConfigs)
		if err != nil {
			fmt.Printf("%d. Error for %s: %s\n", i, name, err)
			continue
		}

		err = cmdutil.WriteFile(name, []byte(processed))
		if err != nil {
			fmt.Printf("%d. Error for %s: %s\n", i, name, err)
			continue
		}
		fmt.Printf("%d. Success for %s", i, name)
	}

	return nil
}

func (m *macroCommand) processWithTemplate(name string, tctx map[string]any) (string, error) {
	query, err := cmdutil.ReadFile(name, nil)
	if err != nil {
		return "", err
	}

	qry, err := m.compiler.CompileString(string(query), tctx)
	if err != nil {
		return "", fmt.Errorf("error in compiling query for file %s: %v", name, err)
	}

	return qry, nil
}
