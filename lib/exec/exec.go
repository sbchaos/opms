package exec

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
)

func Exec(cmd string, args ...string) (stdOut, stdErr bytes.Buffer, err error) {
	pth, err := Path(cmd)
	if err != nil {
		err = fmt.Errorf("could not find %s executable in PATH. error: %w", cmd, err)
		return
	}
	err = Run(context.Background(), pth, nil, nil, &stdOut, &stdErr, args)
	return
}

// ExecContext invokes a gh command in a subprocess and captures the output and error streams.
func ExecContext(ctx context.Context, cmd string, args ...string) (stdout, stderr bytes.Buffer, err error) {
	cmdExe, err := Path(cmd)
	if err != nil {
		err = fmt.Errorf("could not find %s executable in PATH. error: %w", cmd, err)
		return
	}
	err = Run(ctx, cmdExe, nil, nil, &stdout, &stderr, args)
	return
}

// ExecInteractive invokes a gh command in a subprocess with its stdin, stdout, and stderr streams connected to
// those of the parent process. This is suitable for running a command with interactive prompts.
func ExecInteractive(ctx context.Context, cmd string, args ...string) error {
	cmdExe, err := Path(cmd)
	if err != nil {
		return err
	}
	return Run(ctx, cmdExe, nil, os.Stdin, os.Stdout, os.Stderr, args)
}

func Path(cmd string) (string, error) {
	lookPath, err := exec.LookPath(cmd)
	if errors.Is(err, exec.ErrDot) {
		return lookPath, nil
	}
	return lookPath, err
}

func Run(ctx context.Context, exe string, env []string, stdin io.Reader, stdout, stderr io.Writer, args []string) error {
	cmd := exec.CommandContext(ctx, exe, args...)
	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	if env != nil {
		cmd.Env = env
	}
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("execution failed: %w", err)
	}
	return nil
}
