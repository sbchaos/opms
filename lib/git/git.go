package git

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
)

func Exec(args ...string) (stdOut, stdErr bytes.Buffer, err error) {
	pth, err := path()
	if err != nil {
		err = fmt.Errorf("could not find git executable in PATH. error: %w", err)
		return
	}
	return run(pth, nil, args...)
}

func path() (string, error) {
	lookPath, err := exec.LookPath("git")
	if errors.Is(err, exec.ErrDot) {
		return lookPath, nil
	}
	return lookPath, err
}

func run(path string, env []string, args ...string) (stdOut, stdErr bytes.Buffer, err error) {
	cmd := exec.Command(path, args...)
	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr
	if env != nil {
		cmd.Env = env
	}
	err = cmd.Run()
	if err != nil {
		err = fmt.Errorf("failed to run git: %s. error: %w", stdErr.String(), err)
		return
	}
	return
}
