package term

import (
	"io"
	"os"
	"strings"

	"golang.org/x/term"
)

// Term represents information about the terminal that a process is connected to.
type Term struct {
	in           *os.File
	out          *os.File
	errOut       *os.File
	isTTY        bool
	colorEnabled bool
	is256enabled bool
	hasTrueColor bool
	width        int
	widthPercent int
}

// FromEnv initializes a Term from [os.Stdout] and environment variables:
//   - TERM
//   - COLORTERM
func FromEnv(widthOverride int, widthPercent int) Term {
	var stdoutIsTTY bool
	var isColorEnabled bool

	stdoutIsTTY = IsTerminal(os.Stdout)
	isColorEnabled = !IsColorDisabled() && stdoutIsTTY

	return Term{
		in:           os.Stdin,
		out:          os.Stdout,
		errOut:       os.Stderr,
		isTTY:        stdoutIsTTY,
		colorEnabled: isColorEnabled,
		is256enabled: is256ColorSupported(),
		hasTrueColor: isTrueColorSupported(),
		width:        widthOverride,
		widthPercent: widthPercent,
	}
}

// In is the reader reading from standard input.
func (t Term) In() io.Reader {
	return t.in
}

// Out is the writer writing to standard output.
func (t Term) Out() io.Writer {
	return t.out
}

// ErrOut is the writer writing to standard error.
func (t Term) ErrOut() io.Writer {
	return t.errOut
}

// IsTerminalOutput returns true if standard output is connected to a terminal.
func (t Term) IsTerminalOutput() bool {
	return t.isTTY
}

// IsColorEnabled reports whether it's safe to output ANSI color sequences, depending on IsTerminalOutput
// and environment variables.
func (t Term) IsColorEnabled() bool {
	return t.colorEnabled
}

// Is256ColorSupported reports whether the terminal advertises ANSI 256 color codes.
func (t Term) Is256ColorSupported() bool {
	return t.is256enabled
}

// IsTrueColorSupported reports whether the terminal advertises support for ANSI true color sequences.
func (t Term) IsTrueColorSupported() bool {
	return t.hasTrueColor
}

// Size returns the width and height of the terminal that the current process is attached to.
// In case of errors, the numeric values returned are -1.
func (t Term) Size(def int) (int, int) {
	if t.width > 0 {
		return t.width, -1
	}

	ttyOut := t.out
	if ttyOut == nil || !IsTerminal(ttyOut) {
		return def, -1
	}

	width, height, err := terminalSize(ttyOut)
	if err != nil {
		return def, -1
	}

	if t.widthPercent > 0 {
		return int(float64(width) * float64(t.widthPercent) / 100), height
	}

	return width, height
}

// IsTerminal reports whether a file descriptor is connected to a terminal.
func IsTerminal(f *os.File) bool {
	return term.IsTerminal(int(f.Fd()))
}

func terminalSize(f *os.File) (int, int, error) {
	return term.GetSize(int(f.Fd()))
}

// IsColorDisabled returns true if environment variables NO_COLOR or CLICOLOR prohibit usage of color codes
// in terminal output.
func IsColorDisabled() bool {
	return os.Getenv("NO_COLOR") != "" || os.Getenv("CLICOLOR") == "0"
}

func is256ColorSupported() bool {
	return isTrueColorSupported() ||
		strings.Contains(os.Getenv("TERM"), "256") ||
		strings.Contains(os.Getenv("COLORTERM"), "256")
}

func isTrueColorSupported() bool {
	trm := os.Getenv("TERM")
	colorterm := os.Getenv("COLORTERM")

	return strings.Contains(trm, "24bit") ||
		strings.Contains(trm, "truecolor") ||
		strings.Contains(colorterm, "24bit") ||
		strings.Contains(colorterm, "truecolor")
}
