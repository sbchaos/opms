package cmdutil

import (
	"encoding/json"
	"io"
	"os"
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
