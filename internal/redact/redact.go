package redact

import (
	"fmt"
	"regexp"

	"github.com/yourorg/logslice/internal/parser"
)

// Redactor replaces sensitive patterns in log entry messages with a
// configurable placeholder string.
type Redactor struct {
	patterns    []*regexp.Regexp
	placeholder string
}

// New creates a Redactor that replaces matches of any of the provided regex
// patterns with placeholder. An error is returned if any pattern fails to
// compile or if placeholder is empty.
func New(patterns []string, placeholder string) (*Redactor, error) {
	if placeholder == "" {
		return nil, ErrEmptyPlaceholder
	}
	compiled := make([]*regexp.Regexp, 0, len(patterns))
	for _, p := range patterns {
		re, err := regexp.Compile(p)
		if err != nil {
			return nil, fmt.Errorf("redact: invalid pattern %q: %w", p, err)
		}
		compiled = append(compiled, re)
	}
	return &Redactor{patterns: compiled, placeholder: placeholder}, nil
}

// Apply returns a copy of entry with sensitive substrings replaced. The
// original entry is never mutated.
func (r *Redactor) Apply(entry parser.Entry) parser.Entry {
	out := entry
	for _, re := range r.patterns {
		out.Message = re.ReplaceAllString(out.Message, r.placeholder)
	}
	return out
}

// ApplyAll processes a slice of entries and returns redacted copies.
func (r *Redactor) ApplyAll(entries []parser.Entry) []parser.Entry {
	result := make([]parser.Entry, len(entries))
	for i, e := range entries {
		result[i] = r.Apply(e)
	}
	return result
}
