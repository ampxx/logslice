package output_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/logslice/internal/output"
	"github.com/user/logslice/internal/parser"
)

func ts(s string) *time.Time {
	t, _ := time.Parse(time.RFC3339, s)
	return &t
}

func makeEntry(raw string, t *time.Time) parser.Entry {
	return parser.Entry{Raw: raw, Timestamp: t}
}

func TestFormatter_Default(t *testing.T) {
	var buf bytes.Buffer
	f := output.NewFormatter(&buf, output.FormatDefault)
	e := makeEntry("2024-01-02T15:04:05Z info hello", ts("2024-01-02T15:04:05Z"))
	if err := f.Write(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "info hello") {
		t.Errorf("expected raw line in output, got: %q", buf.String())
	}
}

func TestFormatter_JSON(t *testing.T) {
	var buf bytes.Buffer
	f := output.NewFormatter(&buf, output.FormatJSON)
	e := makeEntry("2024-01-02T15:04:05Z error boom", ts("2024-01-02T15:04:05Z"))
	if err := f.Write(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, `"timestamp"`) || !strings.Contains(out, `"raw"`) {
		t.Errorf("expected JSON keys in output, got: %q", out)
	}
}

func TestFormatter_Timestamp(t *testing.T) {
	var buf bytes.Buffer
	f := output.NewFormatter(&buf, output.FormatTimestamp)
	e := makeEntry("2024-03-01T10:00:00Z warn slow query", ts("2024-03-01T10:00:00Z"))
	if err := f.Write(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.HasPrefix(out, "2024-03-01T10:00:00Z") {
		t.Errorf("expected line to start with RFC3339 timestamp, got: %q", out)
	}
}

func TestFormatter_WriteAll(t *testing.T) {
	var buf bytes.Buffer
	f := output.NewFormatter(&buf, output.FormatDefault)
	entries := []parser.Entry{
		makeEntry("line one", nil),
		makeEntry("line two", nil),
	}
	if err := f.WriteAll(entries); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Errorf("expected 2 lines, got %d", len(lines))
	}
}

func TestFormatter_NilTimestamp_JSON(t *testing.T) {
	var buf bytes.Buffer
	f := output.NewFormatter(&buf, output.FormatJSON)
	e := makeEntry("no timestamp line", nil)
	if err := f.Write(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), `"timestamp":null`) {
		t.Errorf("expected null timestamp in JSON, got: %q", buf.String())
	}
}
