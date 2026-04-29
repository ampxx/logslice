// Package truncate provides line-length truncation for log entries.
// It trims the raw text of each entry to a configurable maximum number of
// runes, appending an ellipsis when truncation occurs.
package truncate

import (
	"errors"
	"unicode/utf8"

	"github.com/dkoosis/logslice/internal/parser"
)

// ErrInvalidMaxLen is returned when maxLen is less than 1.
var ErrInvalidMaxLen = errors.New("truncate: maxLen must be >= 1")

// Truncator shortens entry raw text to at most MaxLen runes.
type Truncator struct {
	maxLen int
	suffix string
}

// New creates a Truncator that limits raw text to maxLen runes.
// suffix is appended whenever text is shortened (e.g. "…").
func New(maxLen int, suffix string) (*Truncator, error) {
	if maxLen < 1 {
		return nil, ErrInvalidMaxLen
	}
	return &Truncator{maxLen: maxLen, suffix: suffix}, nil
}

// Apply returns a new entry whose Raw field is truncated to t.maxLen runes.
// All other fields are preserved unchanged.
func (t *Truncator) Apply(e parser.Entry) parser.Entry {
	if utf8.RuneCountInString(e.Raw) <= t.maxLen {
		return e
	}
	runes := []rune(e.Raw)
	cut := t.maxLen
	// Reserve room for suffix if it fits.
	suffixLen := utf8.RuneCountInString(t.suffix)
	if suffixLen > 0 && cut > suffixLen {
		cut -= suffixLen
	}
	e.Raw = string(runes[:cut]) + t.suffix
	return e
}

// ApplyAll truncates every entry in the slice and returns a new slice.
func (t *Truncator) ApplyAll(entries []parser.Entry) []parser.Entry {
	out := make([]parser.Entry, len(entries))
	for i, e := range entries {
		out[i] = t.Apply(e)
	}
	return out
}
