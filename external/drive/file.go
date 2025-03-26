package drive

import (
	"io"

	"google.golang.org/api/drive/v3"

	"github.com/sbchaos/opms/lib/cmdutil"
)

type File struct {
	DriveFile *drive.File
	Path      string
}

func (f File) IsFolder() bool {
	return f.DriveFile.MimeType == "application/vnd.google-apps.folder"
}

func (f File) ToFolder(check bool, allowed map[string]struct{}, checksum map[string]string) Folder {
	return Folder{
		ID:          f.DriveFile.Id,
		Path:        f.Path,
		CheckExt:    check,
		AllowedExt:  allowed,
		ChecksumMap: checksum,
	}
}

func (f File) Name() string {
	return f.DriveFile.Name
}

func (f File) Extension() string {
	return f.DriveFile.FileExtension
}

func (f File) Download(d *drive.Service) error {
	call := d.Files.Get(f.DriveFile.Id)
	subFile, err := call.Download()
	if err != nil {
		return err
	}
	defer subFile.Body.Close()

	contents, err := io.ReadAll(subFile.Body)
	if err != nil {
		return err
	}

	return cmdutil.WriteFile(f.Path, contents)
}
