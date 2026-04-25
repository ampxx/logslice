// Package highlight provides ANSI colour highlighting for matched
// substrings within log entry lines.
package highlight

import (
	"regexp"
	"strings"
)

// ANSI escape codes.
const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Yellow = "\033[33m"
	Cyan   = "\033[36m"
	Bold   = "\033[1m"
)

// Color wraps s with the given ANSI color code and resets afterwards.
func Color(code, s string) string {
	return code + s + Reset
}

// Highlighter applies colour to substrings of a line that match a
// compiled regular expression.
type Highlighter struct {
	re    *regexp.Regexp
	color string
}

// New returns a Highlighter that marks matches of pattern with color.
// If pattern is empty, New returns a no-op Highlighter.
func New(pattern, color string) (*Highlighter, error) {
	if pattern == "" {
		return &Highlighter{}, nil
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	return &Highlighter{re: re, color: color}, nil
}

// Apply returns line with every match of the pattern wrapped in the
// configured ANSI colour sequence. If no pattern was set it returns
// line unchanged.
func (h *Highlighter) Apply(line string) string {
	if h.re == nil {
		return line
	}
	var sb strings.Builder
	last := 0
	for _, loc := range h.re.FindAllStringIndex(line, -1) {
		sb.WriteString(line[last:loc[0]])
		sb.WriteString(Color(h.color, line[loc[0]:loc[1]]))
		last = loc[1]
	}
	sb.WriteString(line[last:])
	return sb.String()
}
