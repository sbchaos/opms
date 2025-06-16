package drive

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/sbchaos/opms/external/drive"
	"github.com/sbchaos/opms/external/gcp"
	"github.com/sbchaos/opms/lib/config"
)

type deleteCommand struct {
	cfg *config.Config

	provider *gcp.ClientProvider

	folderID string
	fileName string
	proj     string
}

func NewDeleteCommand(cfg *config.Config) *cobra.Command {
	del := &deleteCommand{cfg: cfg}

	cmd := &cobra.Command{
		Use:     "delete",
		Short:   "Delete a file from drive",
		Example: "opms drive delete",
		RunE:    del.RunE,
	}

	cmd.Flags().StringVarP(&del.folderID, "folder-id", "f", "", "Drive folder ID")
	cmd.Flags().StringVarP(&del.proj, "project", "p", "", "Project")
	cmd.Flags().StringVarP(&del.fileName, "file-name", "n", "", "File name")

	return cmd
}

func (r *deleteCommand) RunE(_ *cobra.Command, _ []string) error {
	pr, err := gcp.NewClientProvider(r.cfg)
	if err != nil {
		return err
	}
	r.provider = pr

	if r.folderID == "" {
		return errors.New("--folder-id is required")
	}

	root := drive.Folder{
		ID:       r.folderID,
		CheckExt: false,
	}

	service, err := pr.GetDriveClient(r.proj)
	if err != nil {
		return err
	}

	files, err := root.List(service)
	if err != nil {
		return err
	}

	for _, f1 := range files {
		if f1.Name() == r.fileName {
			err = f1.Delete(service)
			if err != nil {
				return err
			}
			fmt.Printf("Deleted file %s\n", f1.Name())
		}
	}

	return nil
}
