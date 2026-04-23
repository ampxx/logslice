package parser

import (
	"bufio"
	"io"
	"regexp"
	"strings"
)

// timestampRe matches common log timestamp prefixes at the start of a line.
var timestampRe = regexp.MustCompile(
	`^(\d{4}[-/]\d{2}[-/]\d{2}[T ]\d{2}:\d{2}:\d{2}(?:Z|[+-]\d{2}:\d{2})?|\w{3}\s+\d+\s+\d{2}:\d{2}:\d{2})`,
)

// levelRe captures common log level keywords.
var levelRe = regexp.MustCompile(`(?i)\b(DEBUG|INFO|WARN(?:ING)?|ERROR|FATAL|TRACE)\b`)

// LineParser parses log lines from an io.Reader into Entry values.
type LineParser struct {
	reader io.Reader
}

// NewLineParser creates a new LineParser reading from r.
func NewLineParser(r io.Reader) *LineParser {
	return &LineParser{reader: r}
}

// Parse reads all lines from the underlying reader and returns parsed entries.
// Lines that cannot have a timestamp extracted are still included with a zero
// timestamp so that no log data is silently dropped.
func (p *LineParser) Parse() ([]Entry, error) {
	var entries []Entry
	scanner := bufio.NewScanner(p.reader)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}
		entries = append(entries, parseLine(line))
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return entries, nil
}

// parseLine converts a raw log line string into an Entry.
func parseLine(line string) Entry {
	e := Entry{Raw: line}

	if m := timestampRe.FindString(line); m != "" {
		if t, _, err := ParseTimestamp(m); err == nil {
			e.Timestamp = t
		}
	}

	if m := levelRe.FindString(line); m != "" {
		e.Level = strings.ToUpper(m)
	}

	return e
}
