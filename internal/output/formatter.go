// Package output provides formatting and writing utilities for log entries.
package output

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/user/logslice/internal/parser"
)

// Format controls how log entries are rendered.
type Format string

const (
	// FormatDefault preserves the original log line as-is.
	FormatDefault Format = "default"
	// FormatJSON renders each entry as a JSON object.
	FormatJSON Format = "json"
	// FormatTimestamp renders only the timestamp and message.
	FormatTimestamp Format = "timestamp"
)

// Formatter writes log entries to an io.Writer in a chosen format.
type Formatter struct {
	w      io.Writer
	format Format
}

// NewFormatter creates a Formatter that writes to w using the given format.
func NewFormatter(w io.Writer, format Format) *Formatter {
	return &Formatter{w: w, format: format}
}

// Write renders a single entry according to the formatter's format.
func (f *Formatter) Write(e parser.Entry) error {
	var line string
	switch f.format {
	case FormatJSON:
		line = entryToJSON(e)
	case FormatTimestamp:
		ts := ""
		if e.Timestamp != nil {
			ts = e.Timestamp.Format(time.RFC3339)
		}
		line = fmt.Sprintf("%s %s", ts, strings.TrimSpace(e.Raw))
	default:
		line = e.Raw
	}
	_, err := fmt.Fprintln(f.w, line)
	return err
}

// WriteAll renders every entry in the slice.
func (f *Formatter) WriteAll(entries []parser.Entry) error {
	for _, e := range entries {
		if err := f.Write(e); err != nil {
			return err
		}
	}
	return nil
}

func entryToJSON(e parser.Entry) string {
	ts := "null"
	if e.Timestamp != nil {
		ts = fmt.Sprintf("%q", e.Timestamp.Format(time.RFC3339))
	}
	raw := strings.ReplaceAll(e.Raw, `"`, `\"`)
	return fmt.Sprintf(`{"timestamp":%s,"raw":"%s"}`, ts, strings.TrimSpace(raw))
}
