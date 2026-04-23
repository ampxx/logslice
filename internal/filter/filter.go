// Package filter provides time-range and pattern-based filtering
// for parsed log entries.
package filter

import (
	"regexp"
	"time"

	"github.com/yourorg/logslice/internal/parser"
)

// Options holds the criteria used to filter log entries.
type Options struct {
	// From filters out entries before this time (inclusive). Zero value means no lower bound.
	From time.Time
	// To filters out entries after this time (inclusive). Zero value means no upper bound.
	To time.Time
	// Pattern is an optional regular expression that must match the entry's raw line.
	Pattern *regexp.Regexp
}

// Filter evaluates a single log entry against the provided Options.
// It returns true when the entry passes all active criteria.
func Filter(entry parser.Entry, opts Options) bool {
	if !opts.From.IsZero() && entry.Timestamp.Before(opts.From) {
		return false
	}
	if !opts.To.IsZero() && entry.Timestamp.After(opts.To) {
		return false
	}
	if opts.Pattern != nil && !opts.Pattern.MatchString(entry.Raw) {
		return false
	}
	return true
}

// FilterAll applies Filter to every entry in the slice and returns those that match.
func FilterAll(entries []parser.Entry, opts Options) []parser.Entry {
	result := make([]parser.Entry, 0, len(entries))
	for _, e := range entries {
		if Filter(e, opts) {
			result = append(result, e)
		}
	}
	return result
}
