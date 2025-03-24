package drive

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"google.golang.org/api/drive/v3"

	"github.com/sbchaos/opms/lib/cmdutil"
	"github.com/sbchaos/opms/lib/pool"
)

type Folder struct {
	ID   string
	Path string

	CheckExt   bool
	AllowedExt map[string]struct{}

	ChecksumMap map[string]string
}

func (f Folder) CanDownload(ext string) bool {
	if !f.CheckExt {
		return true
	}

	_, ok := f.AllowedExt[ext]
	return ok
}

func (f Folder) List(d *drive.Service) ([]File, error) {
	children := make([]File, 0)
	pageToken := ""
	for {
		// TODO: This will  also fetch the folders, we are not getting subfolders
		q := fmt.Sprintf("'%s' in parents", f.ID)
		req := d.Files.
			List().
			IncludeItemsFromAllDrives(true).
			SupportsAllDrives(true).
			Fields("nextPageToken, files(id, name, mimeType, fileExtension, modifiedTime, sha1Checksum)").
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

		for _, f1 := range r.Files {
			newPath := filepath.Join(f.Path, f1.Name)
			children = append(children, File{
				DriveFile: f1,
				Path:      newPath,
			})
		}

		pageToken = r.NextPageToken
		if pageToken == "" {
			break
		}
	}

	return children, nil
}

func (f Folder) Download(d *drive.Service, jobs chan pool.Job[string]) error {
	children, err := f.List(d)
	if err != nil {
		return err
	}

	for _, child := range children {
		if child.IsFolder() {
			fold := child.ToFolder(f.CheckExt, f.AllowedExt, nil)
			err = fold.Download(d, jobs)
			if err != nil {
				fmt.Printf("An error occurred: %v\n", err)
			}

			continue
		}

		if f.CanDownload(child.Extension()) {
			jobs <- func() pool.JobResult[string] {
				err = child.Download(d)
				return pool.JobResult[string]{
					Output: child.Name(),
					Err:    err,
				}
			}
		}
	}
	return nil
}

func (f Folder) Sync(d *drive.Service, jobs chan pool.Job[string]) error {
	children, err := f.List(d)
	if err != nil {
		return err
	}

	localFileMap := map[string]string{}
	files, err := cmdutil.ListFiles(f.Path)
	if err == nil {
		for _, localFile := range files {
			localFileMap[localFile] = ""
		}
	}

	for _, child := range children {
		if _, ok := localFileMap[child.Path]; ok {
			delete(localFileMap, child.Path)
		}

		if child.IsFolder() {
			os.MkdirAll(child.Path, fs.ModePerm)
			fold := child.ToFolder(f.CheckExt, f.AllowedExt, f.ChecksumMap)
			err = fold.Sync(d, jobs)
			if err != nil {
				fmt.Printf("An error occurred: %v\n", err)
			}

			continue
		}

		if f.CanDownload(child.Extension()) {
			download := false

			oldCheckSum, ok := f.ChecksumMap[child.Path]
			if !ok || oldCheckSum != child.DriveFile.Sha1Checksum {
				download = true
			}

			if download {
				fmt.Printf("Download %s\n", child.Path)
				f.ChecksumMap[child.Path] = child.DriveFile.Sha1Checksum
				jobs <- func() pool.JobResult[string] {
					err = child.Download(d)
					return pool.JobResult[string]{
						Output: child.Name(),
						Err:    err,
					}
				}
			}
		}
	}

	for name, _ := range localFileMap {
		fmt.Printf("Removing %s\n", name)
		os.RemoveAll(name)
		delete(f.ChecksumMap, name)
	}

	return nil
}
