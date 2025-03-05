package cmdutil

import (
	"os"
	"path/filepath"
)

func WriteFileAndDir(filename string, data []byte) error {
	dir := filepath.Dir(filename)
	err := CreateDir(dir)
	if err != nil {
		return err
	}

	return WriteFile(filename, data)
}

func CreateDir(dirpath string) error {
	return os.MkdirAll(dirpath, os.ModePerm)
}

func WriteFile(filename string, data []byte) error {
	return os.WriteFile(filename, data, 0666)
}
