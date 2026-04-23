// Package parser provides log entry parsing functionality for logslice.
package parser

import (
	"fmt"
	"time"
)

// Entry represents a single parsed log line with its timestamp and raw content.
type Entry struct {
	Timestamp time.Time
	Raw       string
	Level     string
}

// String returns a human-readable representation of the log entry.
func (e Entry) String() string {
	return fmt.Sprintf("[%s] %s: %s", e.Timestamp.Format(time.RFC3339), e.Level, e.Raw)
}

// CommonLayouts holds timestamp formats commonly found in log files.
var CommonLayouts = []string{
	"2006-01-02T15:04:05Z07:00",
	"2006-01-02T15:04:05",
	"2006-01-02 15:04:05",
	"02/Jan/2006:15:04:05 -0700",
	"Jan 2 15:04:05",
	"2006/01/02 15:04:05",
}

// ParseTimestamp attempts to parse a timestamp string using known layouts.
// Returns the parsed time and the matched layout, or an error if none match.
func ParseTimestamp(s string) (time.Time, string, error) {
	for _, layout := range CommonLayouts {
		t, err := time.Parse(layout, s)
		if err == nil {
			return t, layout, nil
		}
	}
	return time.Time{}, "", fmt.Errorf("unable to parse timestamp: %q", s)
}
