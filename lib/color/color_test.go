package color_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sbchaos/opms/lib/color"
)

func TestColorFromRGB(t *testing.T) {
	tests := []struct {
		name  string
		hex   string
		text  string
		wants string
		cs    *color.Scheme
	}{
		{
			name:  "truecolor",
			hex:   "fc0303",
			text:  "red",
			wants: "\033[38;2;252;3;3mred\033[0m",
			cs:    color.NewColorScheme(true, true, true),
		},
		{
			name:  "no truecolor",
			hex:   "fc0303",
			text:  "red",
			wants: "red",
			cs:    color.NewColorScheme(true, true, false),
		},
		{
			name:  "no color",
			hex:   "fc0303",
			text:  "red",
			wants: "red",
			cs:    color.NewColorScheme(false, false, false),
		},
		{
			name:  "invalid hex",
			hex:   "fc0",
			text:  "red",
			wants: "red",
			cs:    color.NewColorScheme(false, false, false),
		},
	}

	for _, tt := range tests {
		fn := tt.cs.ColorFromRGB(tt.hex)
		assert.Equal(t, tt.wants, fn(tt.text))
	}
}

func TestHexToRGB(t *testing.T) {
	tests := []struct {
		name  string
		hex   string
		text  string
		wants string
		cs    *color.Scheme
	}{
		{
			name:  "truecolor",
			hex:   "fc0303",
			text:  "red",
			wants: "\033[38;2;252;3;3mred\033[0m",
			cs:    color.NewColorScheme(true, true, true),
		},
		{
			name:  "no truecolor",
			hex:   "fc0303",
			text:  "red",
			wants: "red",
			cs:    color.NewColorScheme(true, true, false),
		},
		{
			name:  "no color",
			hex:   "fc0303",
			text:  "red",
			wants: "red",
			cs:    color.NewColorScheme(false, false, false),
		},
		{
			name:  "invalid hex",
			hex:   "fc0",
			text:  "red",
			wants: "red",
			cs:    color.NewColorScheme(false, false, false),
		},
	}

	for _, tt := range tests {
		output := tt.cs.HexToRGB(tt.hex, tt.text)
		assert.Equal(t, tt.wants, output)
	}
}

func TestColors(t *testing.T) {
	schm := color.NewColorScheme(true, true, true)
	for _, c := range []int{color.Black, color.Red, color.Green, color.Yellow, color.Blue, color.Magenta, color.Cyan, color.LightGray, color.DarkGray, color.LightRed, color.LightGreen, color.LightYellow, color.LightBlue, color.LightMagenta, color.LightCyan, color.White} {
		colored := schm.Colorize(c, color.Normal, "text")
		fmt.Println(colored)
	}

	//fmt.Println(schm.Colorize(255, color.HighlightStyle, "text"))
}
