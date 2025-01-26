// Package text is a set of utility functions for text processing and outputting to the terminal.
package text

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

const (
	ellipsis            = "..."
	minWidthForEllipsis = len(ellipsis) + 2
)

var indentRE = regexp.MustCompile(`(?m)^`)

// Indent returns a copy of the string s with indent prefixed to it, will apply indent
// to each line of the string.
func Indent(s, indent string) string {
	if len(strings.TrimSpace(s)) == 0 {
		return s
	}
	return indentRE.ReplaceAllLiteralString(s, indent)
}

// TruncateMultiline returns a copy of the string s that has been shortened to fit the maximum
// display width. If string s has multiple lines the first line will be shortened and all others
// removed.
func TruncateMultiline(maxWidth int, s string) string {
	if i := strings.IndexAny(s, "\r\n"); i >= 0 {
		s = s[:i] + ellipsis
	}
	return Truncate(maxWidth, s)
}

// Truncate returns a copy of the string s that has been shortened to fit the maximum display width.
func Truncate(maxWidth int, s string) string {
	w := len(s)
	if w <= maxWidth {
		return s
	}
	tail := ""
	if maxWidth >= minWidthForEllipsis {
		tail = ellipsis
	}
	r := s[0:uint(maxWidth-3)] + tail
	if len(r) < maxWidth {
		r += " "
	}
	return r
}

// PadRight returns a copy of the string s that has been padded on the right with whitespace to fit
// the maximum display width.
func PadRight(maxWidth int, s string) string {
	if padWidth := maxWidth - len(s); padWidth > 0 {
		s += strings.Repeat(" ", padWidth)
	}
	return s
}

// Pluralize returns a concatenated string with num and the plural form of thing if necessary.
func Pluralize(num int, thing string) string {
	if num == 1 {
		return fmt.Sprintf("%d %s", num, thing)
	}
	return fmt.Sprintf("%d %ss", num, thing)
}

func fmtDuration(amount int, unit string) string {
	return fmt.Sprintf("about %s ago", Pluralize(amount, unit))
}

// RelativeTimeAgo returns a human-readable string of the time duration between a and b that is estimated
// to the nearest unit of time.
func RelativeTimeAgo(a, b time.Time) string {
	ago := a.Sub(b)

	return TimeAgo(ago)
}

func TimeAgo(ago time.Duration) string {
	if ago < time.Minute {
		return "less than a minute ago"
	}
	if ago < time.Hour {
		return fmtDuration(int(ago.Minutes()), "minute")
	}
	if ago < 24*time.Hour {
		return fmtDuration(int(ago.Hours()), "hour")
	}
	if ago < 30*24*time.Hour {
		return fmtDuration(int(ago.Hours())/24, "day")
	}
	if ago < 365*24*time.Hour {
		return fmtDuration(int(ago.Hours())/24/30, "month")
	}

	return fmtDuration(int(ago.Hours()/24/365), "year")
}
