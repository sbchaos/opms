package table_test

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/sbchaos/opms/lib/printers/table"
)

func E1TablePrinter() {
	// information about the terminal can be obtained using the [pkg/term] package
	isTTY := true
	termWidth := 14
	red := func(s string) string {
		return "\x1b[31m" + s + "\x1b[m"
	}

	t := table.New(os.Stdout, isTTY, termWidth)
	t.AddField("9", table.WithTruncate(nil))
	t.AddField("hello")
	t.EndRow()
	t.AddField("10", table.WithTruncate(nil))
	t.AddField("long description", table.WithColor(red))
	t.EndRow()
	if err := t.Render(); err != nil {
		log.Fatal(err)
	}
	// stdout now contains:
	// 9   hello
	// 10  long de...
}

func Test_ttyTablePrinter_autoTruncate(t *testing.T) {
	buf := bytes.Buffer{}
	tp := table.New(&buf, true, 5)

	tp.AddField("1")
	tp.AddField("hello")
	tp.EndRow()
	tp.AddField("2")
	tp.AddField("world")
	tp.EndRow()

	err := tp.Render()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "1  he\n2  wo\n"
	if buf.String() != expected {
		t.Errorf("expected: %q, got: %q", expected, buf.String())
	}
}

func Test_ttyTablePrinter_WithTruncate(t *testing.T) {
	buf := bytes.Buffer{}
	tp := table.New(&buf, true, 15)

	tp.AddField("long SHA", table.WithTruncate(nil))
	tp.AddField("hello")
	tp.EndRow()
	tp.AddField("another SHA", table.WithTruncate(nil))
	tp.AddField("world")
	tp.EndRow()

	err := tp.Render()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "long SHA     he\nanother SHA  wo\n"
	if buf.String() != expected {
		t.Errorf("expected: %q, got: %q", expected, buf.String())
	}
}

func Test_ttyTablePrinter_AddHeader(t *testing.T) {
	buf := bytes.Buffer{}
	tp := table.New(&buf, true, 80)

	tp.AddHeader([]string{"ONE", "TWO", "THREE"}, table.WithColor(func(s string) string {
		return fmt.Sprintf("\x1b[4m%s\x1b[m", s)
	}))
	// Subsequent calls to AddHeader are ignored.
	tp.AddHeader([]string{"SHOULD", "NOT", "EXIST"})

	tp.AddField("hello")
	tp.AddField("beautiful")
	tp.AddField("people")
	tp.EndRow()

	err := tp.Render()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := `
		\x1b[4mONE  \x1b[m  \x1b[4mTWO      \x1b[m  \x1b[4mTHREE\x1b[m
		hello  beautiful  people
	`
	if buf.String() != expected {
		t.Errorf("expected: %q, got: %q", expected, buf.String())
	}
}

func Test_ttyTablePrinter_WithPadding(t *testing.T) {
	buf := bytes.Buffer{}
	tp := table.New(&buf, true, 80)

	// Center the headers.
	tp.AddHeader([]string{"A", "B", "C"}, table.WithPadding(func(width int, s string) string {
		left := (width - len(s)) / 2
		return strings.Repeat(" ", left) + s + strings.Repeat(" ", width-left-len(s))
	}))

	tp.AddField("hello")
	tp.AddField("beautiful")
	tp.AddField("people")
	tp.EndRow()

	err := tp.Render()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := `
		  A        B        C   
		hello  beautiful  people
	`
	if buf.String() != expected {
		t.Errorf("expected: %q, got: %q", expected, buf.String())
	}
}

func Test_tsvTablePrinter(t *testing.T) {
	buf := bytes.Buffer{}
	tp := table.New(&buf, false, 0)

	tp.AddField("1")
	tp.AddField("hello")
	tp.EndRow()
	tp.AddField("2")
	tp.AddField("world")
	tp.EndRow()

	err := tp.Render()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "1\thello\n2\tworld\n"
	if buf.String() != expected {
		t.Errorf("expected: %q, got: %q", expected, buf.String())
	}
}

func Test_tsvTablePrinter_AddHeader(t *testing.T) {
	buf := bytes.Buffer{}
	tp := table.New(&buf, false, 0)

	// Headers are not output in TSV output.
	tp.AddHeader([]string{"ONE", "TWO", "THREE"})

	tp.AddField("hello")
	tp.AddField("beautiful")
	tp.AddField("people")
	tp.EndRow()
	tp.AddField("1")
	tp.AddField("2")
	tp.AddField("3")
	tp.EndRow()

	err := tp.Render()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "hello\tbeautiful\tpeople\n1\t2\t3\n"
	if buf.String() != expected {
		t.Errorf("expected: %q, got: %q", expected, buf.String())
	}
}
