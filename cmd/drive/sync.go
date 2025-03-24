package drive

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/sbchaos/opms/external/drive"
	"github.com/sbchaos/opms/external/gcp"
	"github.com/sbchaos/opms/lib/cmdutil"
	"github.com/sbchaos/opms/lib/config"
	"github.com/sbchaos/opms/lib/pool"
)

type syncCommand struct {
	cfg *config.Config

	provider *gcp.ClientProvider

	folderID string

	fileExt string
	output  string
	proj    string
	workers int

	checksumPath string
}

func NewSyncCommand(cfg *config.Config) *cobra.Command {
	sync := &syncCommand{cfg: cfg}

	cmd := &cobra.Command{
		Use:     "sync",
		Short:   "Sync content of a directory from drive",
		Example: "opms drive sync",
		RunE:    sync.RunE,
	}

	cmd.Flags().StringVarP(&sync.folderID, "folder-id", "f", "", "Drive folder ID")
	cmd.Flags().StringVarP(&sync.proj, "project", "p", "", "Project")
	cmd.Flags().StringVarP(&sync.fileExt, "file-ext", "e", "", "File extensions, comma separated")
	cmd.Flags().StringVarP(&sync.output, "output", "o", "", "Output folder")
	cmd.Flags().StringVarP(&sync.checksumPath, "checksum", "c", "", "Checksum json file")
	cmd.Flags().IntVarP(&sync.workers, "workers", "w", 1, "Number of parallel workers")

	return cmd
}

func (r *syncCommand) RunE(_ *cobra.Command, _ []string) error {
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

	checksum := map[string]string{}
	if r.checksumPath != "" {
		err = cmdutil.ReadJsonFile(r.checksumPath, os.Stdin, &checksum)
		if err != nil {
			return err
		}
	}
	root.ChecksumMap = checksum

	service, err := pr.GetDriveClient(r.proj)
	if err != nil {
		return err
	}

	jobs := make(chan pool.Job[string], 20)
	outChan := pool.StartPool(r.workers, jobs)

	go func() {
		err = root.Sync(service, jobs)
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

	bytes, err := json.Marshal(checksum)
	if err != nil {
		fmt.Printf("Error in marshaling")
	}

	outputJson := filepath.Join(r.output, "checksum.json")
	if r.checksumPath != "" {
		outputJson = r.checksumPath
	}
	cmdutil.WriteFile(outputJson, bytes)
	return nil
}
