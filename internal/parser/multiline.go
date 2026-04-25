package parser

import (
	"strings"
	"time"
)

// MultiLineParser groups continuation lines (lines without a recognisable
// timestamp) with the preceding log entry.  This is common in Java stack
// traces, Python tracebacks, and any logger that wraps long messages across
// several physical lines.
//
// Usage:
//
//	mlp := NewMultiLineParser()
//	for _, raw := range rawLines {
//		if entry, ok := mlp.Feed(raw); ok {
//			// entry is complete and ready to process
//		}
//	}
//	if entry, ok := mlp.Flush(); ok {
//		// handle the last buffered entry
//	}
type MultiLineParser struct {
	lp      *LineParser
	pending *Entry
}

// NewMultiLineParser returns a MultiLineParser backed by a default LineParser.
func NewMultiLineParser() *MultiLineParser {
	return &MultiLineParser{lp: NewLineParser()}
}

// Feed accepts a single raw log line.  When a new timestamped line is
// encountered the previously buffered entry (if any) is returned together with
// true.  Continuation lines are appended to the current buffer.  When no
// complete entry is ready yet the function returns (nil, false).
func (m *MultiLineParser) Feed(line string) (*Entry, bool) {
	// Blank lines are treated as continuations so that empty lines inside a
	// stack trace are preserved.
	if strings.TrimSpace(line) == "" {
		if m.pending != nil {
			m.pending.Raw += "\n" + line
		}
		return nil, false
	}

	next := m.lp.Parse(line)

	// If the line has no timestamp it is a continuation of the current entry.
	if next == nil || next.Timestamp.IsZero() {
		if m.pending != nil {
			m.pending.Raw += "\n" + line
		}
		return nil, false
	}

	// We have a new timestamped line — flush the previous entry first.
	prev := m.pending
	m.pending = next

	if prev != nil {
		return prev, true
	}
	return nil, false
}

// Flush returns any buffered entry that has not yet been emitted.  Call this
// after the input stream is exhausted to ensure the final entry is not lost.
func (m *MultiLineParser) Flush() (*Entry, bool) {
	if m.pending == nil {
		return nil, false
	}
	e := m.pending
	m.pending = nil
	return e, true
}

// Reset clears internal state so the parser can be reused for a new stream.
func (m *MultiLineParser) Reset() {
	m.pending = nil
}

// ParseAll is a convenience helper that processes every line in the provided
// slice and returns all assembled entries.  It is primarily useful in tests
// and small scripts where streaming is not required.
func ParseAll(lines []string) []*Entry {
	mp := NewMultiLineParser()
	var entries []*Entry

	for _, l := range lines {
		if e, ok := mp.Feed(l); ok {
			entries = append(entries, e)
		}
	}
	if e, ok := mp.Flush(); ok {
		entries = append(entries, e)
	}
	return entries
}

// stubTimestamp is used internally when synthesising entries for lines that
// arrive before any timestamped line has been seen.  The zero value of
// time.Time signals "unknown" to downstream filters.
var stubTimestamp = time.Time{}
