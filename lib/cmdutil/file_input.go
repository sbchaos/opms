package cmdutil

import (
	"encoding/json"
	"io"
	"os"
	"path"
	"strings"
)

func ReadFile(filename string, stdin io.ReadCloser) ([]byte, error) {
	if filename == "-" {
		b, err := io.ReadAll(stdin)
		_ = stdin.Close()
		return b, err
	}

	return os.ReadFile(filename)
}

func ReadJsonFile(filename string, stdin io.ReadCloser, v any) error {
	file, err := ReadFile(filename, stdin)
	if err != nil {
		return err
	}

	return json.Unmarshal(file, v)
}

func ReadLines(filename string, stdin io.ReadCloser) ([]string, error) {
	file, err := ReadFile(filename, stdin)
	if err != nil {
		return nil, err
	}

	return strings.Fields(string(file)), nil
}

func ListFiles(dirName string) ([]string, error) {
	return recursiveListFiles(dirName, 1)
}

func ListAllFiles(dirName string) ([]string, error) {
	return recursiveListFiles(dirName, 100)
}

func recursiveListFiles(dirName string, level int) ([]string, error) {
	entries, err := os.ReadDir(dirName)
	if err != nil {
		return nil, err
	}

	var names []string
	for _, entry := range entries {
		if !entry.IsDir() {
			filepath := path.Join(dirName, entry.Name())
			names = append(names, filepath)
			continue
		}

		if level <= 1 {
			continue
		}

		newPath := path.Join(dirName, entry.Name())
		newLevel := level - 1
		n1, errChild := recursiveListFiles(newPath, newLevel)
		names = append(names, n1...)
		if errChild != nil {
			return names, errChild
		}
	}

	return names, nil
}
