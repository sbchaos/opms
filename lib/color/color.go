package color

import (
	"bytes"
	"fmt"
	"strconv"
)

const (
	Start         = "\033["
	Normal        = "0;"
	Bold          = "1;"
	Dim           = "2;"
	Underline     = "4;"
	Blink         = "5;"
	Inverse       = "7;"
	Strikethrough = "9;"

	// DefaultBG is the default background
	DefaultBG = "\033[49m"
	// DefaultFG is the default foreground
	DefaultFG = "\033[39m"
)

const (
	Reset = "\033[0m"

	Magenta      = 35
	LightGray    = 37
	DarkGray     = 90
	LightRed     = 91
	LightGreen   = 92
	LightYellow  = 93
	LightBlue    = 94
	LightMagenta = 95
	LightCyan    = 96

	Black  = 30
	Cyan   = 36
	Red    = 31
	Green  = 32
	Yellow = 33
	Blue   = 34
	White  = 38
)

func gray256(t string) string {
	return fmt.Sprintf("\x1b[%d;5;%dm%s\x1b[m", 38, 242, t)
}

func NewColorScheme(enabled, is256enabled bool, trueColor bool) *Scheme {
	return &Scheme{
		enabled:      enabled,
		is256enabled: is256enabled,
		hasTrueColor: trueColor,
	}
}

type Scheme struct {
	enabled      bool
	is256enabled bool
	hasTrueColor bool
}

func (c *Scheme) Enabled() bool {
	return c.enabled
}

func (c *Scheme) Colorize(color int, style string, t string) string {
	buf := bytes.NewBufferString(Start)
	if style != "" {
		buf.WriteString(style)
	}

	fmt.Fprintf(buf, "38;5;%dm", color)
	buf.WriteString(t)
	buf.WriteString(Reset)
	return buf.String()
}

// ColorFromRGB returns a function suitable for TablePrinter.AddField
// that calls HexToRGB, coloring text if supported by the terminal.
func (c *Scheme) ColorFromRGB(hex string) func(string) string {
	return func(s string) string {
		return c.HexToRGB(hex, s)
	}
}

// HexToRGB uses the given hex to color x if supported by the terminal.
func (c *Scheme) HexToRGB(hex string, x string) string {
	if !c.enabled || !c.hasTrueColor || len(hex) != 6 {
		return x
	}

	r, _ := strconv.ParseInt(hex[0:2], 16, 64)
	g, _ := strconv.ParseInt(hex[2:4], 16, 64)
	b, _ := strconv.ParseInt(hex[4:6], 16, 64)
	return fmt.Sprintf("\033[38;2;%d;%d;%dm%s\033[0m", r, g, b, x)
}
