package gsheet

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/sbchaos/opms/external/gcp"
	"github.com/sbchaos/opms/external/gsheet"
	"github.com/sbchaos/opms/lib/config"
	"github.com/sbchaos/opms/lib/table"
	"github.com/sbchaos/opms/lib/term"
)

type readCommand struct {
	cfg *config.Config

	provider *gcp.ClientProvider

	sheetURL string
	sheetID  string

	sheetRange string

	noFormat bool
	proj     string
}

func NewReadCommand(cfg *config.Config) *cobra.Command {
	read := &readCommand{cfg: cfg}

	cmd := &cobra.Command{
		Use:     "read",
		Short:   "Read content from google sheet",
		Example: "opms gsheets read",
		RunE:    read.RunE,
	}

	cmd.Flags().StringVarP(&read.sheetID, "sheet-id", "i", "", "Sheet ID")
	cmd.Flags().StringVarP(&read.sheetURL, "sheet", "s", "", "Google sheet URL")
	cmd.Flags().StringVarP(&read.sheetRange, "sheetRange", "r", "", "Sheet range")
	cmd.Flags().StringVarP(&read.proj, "project", "p", "", "Project")
	cmd.Flags().BoolVarP(&read.noFormat, "no-format", "f", false, "Do not get formatted data")

	return cmd
}

func (r *readCommand) RunE(_ *cobra.Command, _ []string) error {
	pr, err := gcp.NewClientProvider(r.cfg)
	if err != nil {
		return err
	}
	r.provider = pr

	if r.sheetID == "" && r.sheetURL == "" {
		return errors.New("either --sheet or --sheet-id flag must be set")
	}

	sheetID := r.sheetID
	if sheetID == "" {
		info, err := gsheet.FromURL(r.sheetURL)
		if err != nil {
			return err
		}
		sheetID = info.SheetID
	}

	client, err := pr.GetSheetsClient(r.proj)
	if err != nil {
		return err
	}

	t := term.FromEnv(0, 0)
	size, _ := t.Size(120)

	printer := table.New(os.Stdout, t.IsTerminalOutput(), size)

	content, err := gsheet.GetSheetContent(client, sheetID, r.sheetRange)
	if err != nil {
		return err
	}

	headers := []string{"num"}
	for _, r1 := range content[0] {
		headers = append(headers, fmt.Sprintf("%v", r1))
	}

	printer.AddHeader(headers)
	for i, row := range content[1:] {
		printer.AddField(strconv.Itoa(i + 1))
		for _, r1 := range row {
			printer.AddField(fmt.Sprintf("%s", r1))
		}
		printer.EndRow()
	}

	return printer.Render()
}
