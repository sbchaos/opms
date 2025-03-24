package drive

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/sbchaos/opms/external/drive"
	"github.com/sbchaos/opms/external/gcp"
	"github.com/sbchaos/opms/lib/config"
	"github.com/sbchaos/opms/lib/pool"
)

type downloadCommand struct {
	cfg *config.Config

	provider *gcp.ClientProvider

	folderID string

	fileExt string
	output  string
	proj    string
	workers int
}

func NewDownloadCommand(cfg *config.Config) *cobra.Command {
	download := &downloadCommand{cfg: cfg}

	cmd := &cobra.Command{
		Use:     "download",
		Short:   "Download content from drive",
		Example: "opms drive download",
		RunE:    download.RunE,
	}

	cmd.Flags().StringVarP(&download.folderID, "folder-id", "f", "", "Drive folder ID")
	cmd.Flags().StringVarP(&download.proj, "project", "p", "", "Project")
	cmd.Flags().StringVarP(&download.fileExt, "file-ext", "e", "", "File extensions, comma separated")
	cmd.Flags().StringVarP(&download.output, "output", "o", "", "Output folder")
	cmd.Flags().IntVarP(&download.workers, "workers", "w", 1, "Number of parallel workers")

	return cmd
}

func (r *downloadCommand) RunE(_ *cobra.Command, _ []string) error {
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
		Path:     r.output,
		CheckExt: false,
	}

	allowedExt := make(map[string]struct{})
	if r.fileExt != "" {
		root.CheckExt = true
		parts := strings.SplitSeq(r.fileExt, ",")
		for part := range parts {
			key := strings.TrimSpace(part)
			allowedExt[key] = struct{}{}
		}
		root.AllowedExt = allowedExt
	}

	service, err := pr.GetDriveClient(r.proj)
	if err != nil {
		return err
	}

	jobs := make(chan pool.Job[string], 20)
	outChan := pool.StartPool(r.workers, jobs)

	go func() {
		err = root.Download(service, jobs)
		if err != nil {
			fmt.Println("Error:", err)
		}
		close(jobs)
	}()

	for out := range outChan {
		if out.Err != nil {
			fmt.Fprintf(os.Stderr, "Name: %s, Err: %s\n", out.Output, out.Err)
		}
	}
	return nil
}
