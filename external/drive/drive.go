package drive

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"google.golang.org/api/drive/v3"

	"github.com/sbchaos/opms/lib/cmdutil"
	"github.com/sbchaos/opms/lib/pool"
)

func ListFiles(d *drive.Service, folderID string) ([]*drive.File, error) {
	children := make([]*drive.File, 0)
	pageToken := ""
	for {
		// TODO: This will  also fetch the folders, we are not getting subfolders
		q := fmt.Sprintf("'%s' in parents", folderID)
		req := d.Files.
			List().
			IncludeItemsFromAllDrives(true).
			SupportsAllDrives(true).
			Fields("nextPageToken, files(id, name, mimeType)").
			Q(q)

		// If we have a pageToken set, apply it to the query
		if pageToken != "" {
			req = req.PageToken(pageToken)
		}
		r, err := req.Do()
		if err != nil {
			fmt.Printf("An error occurred: %v\n", err)
			return children, err
		}

		children = append(children, r.Files...)
		pageToken = r.NextPageToken
		if pageToken == "" {
			break
		}
	}

	fmt.Printf("Found %d files\n", len(children))
	return children, nil
}

func DownloadFile(d *drive.Service, file *drive.File, filepath string) error {
	call := d.Files.Get(file.Id)
	subFile, err := call.Download()
	if err != nil {
		return err
	}
	defer subFile.Body.Close()

	contents, err := io.ReadAll(subFile.Body)
	if err != nil {
		return err
	}

	return cmdutil.WriteFileAndDir(filepath, contents)
}

func DownloadFolder(d *drive.Service, folderID string, dirPath string, jobs chan pool.Job[string]) error {
	children, err := ListFiles(d, folderID)
	if err != nil {
		return err
	}

	for _, child := range children {
		if strings.HasPrefix(child.Name, ".") {
			continue
		}

		newPath := filepath.Join(dirPath, child.Name)
		if child.MimeType == "application/vnd.google-apps.folder" {
			err = DownloadFolder(d, child.Id, newPath, jobs)
			if err != nil {
				fmt.Printf("An error occurred: %v\n", err)
			}

			continue
		}

		jobs <- func() pool.JobResult[string] {
			f1 := child
			n1 := newPath
			err = DownloadFile(d, f1, n1)
			return pool.JobResult[string]{
				Output: f1.Name,
				Err:    err,
			}
		}
	}
	return nil
}
