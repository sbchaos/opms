package term_test

import (
	"testing"

	"github.com/sbchaos/opms/lib/term"
)

func TestFromEnv(t *testing.T) {
	tests := []struct {
		name          string
		env           map[string]string
		wantColor     bool
		want256Color  bool
		wantTrueColor bool
	}{
		{
			name: "default",
			env: map[string]string{
				"NO_COLOR":  "",
				"TERM":      "",
				"COLORTERM": "",
			},
			wantColor:     false,
			want256Color:  false,
			wantTrueColor: false,
		},
		{
			name: "no color",
			env: map[string]string{
				"NO_COLOR":  "TRUE",
				"TERM":      "",
				"COLORTERM": "",
			},
			wantColor:     false,
			want256Color:  false,
			wantTrueColor: false,
		},
		{
			name: "has 256-color support",
			env: map[string]string{
				"TERM":      "256-color",
				"COLORTERM": "truecolor",
			},
			wantColor:     false,
			want256Color:  true,
			wantTrueColor: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for key, value := range tt.env {
				t.Setenv(key, value)
			}
			terminal := term.FromEnv(120, 90)
			if got := terminal.IsColorEnabled(); got != tt.wantColor {
				t.Errorf("expected color %v, got %v", tt.wantColor, got)
			}
			if got := terminal.Is256ColorSupported(); got != tt.want256Color {
				t.Errorf("expected 256-color %v, got %v", tt.want256Color, got)
			}
			if got := terminal.IsTrueColorSupported(); got != tt.wantTrueColor {
				t.Errorf("expected truecolor %v, got %v", tt.wantTrueColor, got)
			}
		})
	}
}
