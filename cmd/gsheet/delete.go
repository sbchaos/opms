package gsheet

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/spf13/cobra"

	"github.com/sbchaos/opms/external/gcp"
	"github.com/sbchaos/opms/external/gsheet"
	"github.com/sbchaos/opms/lib/config"
)

type deleteCommand struct {
	cfg *config.Config

	provider *gcp.ClientProvider

	sheetURL string
	sheetID  string

	tabName    string
	tabPattern string

	proj string
}

func NewDeleteCommand(cfg *config.Config) *cobra.Command {
	del := &deleteCommand{cfg: cfg}

	cmd := &cobra.Command{
		Use:     "delete",
		Short:   "Delete google sheet tab",
		Example: "opms gsheets delete",
		RunE:    del.RunE,
	}

	cmd.Flags().StringVarP(&del.sheetID, "sheet-id", "i", "", "Sheet ID")
	cmd.Flags().StringVarP(&del.sheetURL, "sheet", "s", "", "Google sheet URL")
	cmd.Flags().StringVarP(&del.proj, "project", "p", "", "Project")
	cmd.Flags().StringVarP(&del.tabName, "tab-name", "t", "", "tab name to delete")
	cmd.Flags().StringVarP(&del.tabPattern, "tab-pattern", "r", "", "tab name regex")

	return cmd
}

func (d *deleteCommand) RunE(_ *cobra.Command, _ []string) error {
	pr, err := gcp.NewClientProvider(d.cfg)
	if err != nil {
		return err
	}
	d.provider = pr

	if d.sheetID == "" && d.sheetURL == "" {
		return errors.New("either --sheet or --sheet-id flag must be set")
	}

	if d.tabName == "" && d.tabPattern == "" {
		return errors.New("either --tab or --tab-pattern flag must be set")
	}

	var pat *regexp.Regexp = nil
	if d.tabPattern != "" {
		pattern, err := regexp.Compile(d.tabPattern)
		if err != nil {
			return err
		}
		pat = pattern
	}

	sheetID := d.sheetID
	if sheetID == "" {
		info, err := gsheet.FromURL(d.sheetURL)
		if err != nil {
			return err
		}
		sheetID = info.SheetID
	}

	client, err := pr.GetSheetsClient(d.proj)
	if err != nil {
		return err
	}

	sheets, err := gsheet.GetSheets(client, sheetID)
	if err != nil {
		return err
	}

	if d.tabName != "" {
		for _, sheet := range sheets {
			if sheet.Properties.Title == d.tabName {
				err = gsheet.DeleteSheet(client, sheetID, sheet.Properties.SheetId)
				if err == nil {
					fmt.Printf("Deleted tab %s\n", d.tabName)
				}
				return err
			}
		}
	}

	if pat != nil {
		for _, sheet := range sheets {
			title := sheet.Properties.Title
			match := pat.MatchString(title)
			if match {
				err = gsheet.DeleteSheet(client, sheetID, sheet.Properties.SheetId)
				if err != nil {
					fmt.Printf("Error deleting tab: %s: %s\n", title, err)
				} else {
					fmt.Printf("Deleted tab %s\n", title)
				}
			}
		}
	}
	return nil
}
