// Package stats provides aggregation and summary statistics over log entries.
package stats

import (
	"time"

	"github.com/yourorg/logslice/internal/parser"
)

// Summary holds aggregate statistics computed from a set of log entries.
type Summary struct {
	Total      int
	Matched    int
	Skipped    int
	Earliest   time.Time
	Latest     time.Time
	PatternHit int
}

// Collector accumulates log entry statistics.
type Collector struct {
	total      int
	matched    int
	skipped    int
	earliestSet bool
	earliestTS  time.Time
	latestTS    time.Time
	patternHit int
}

// NewCollector returns an initialised Collector.
func NewCollector() *Collector {
	return &Collector{}
}

// Record registers a parsed entry. matched indicates whether the entry passed
// all filters; patternMatch indicates whether it matched the pattern filter.
func (c *Collector) Record(e parser.Entry, matched, patternMatch bool) {
	c.total++
	if matched {
		c.matched++
	} else {
		c.skipped++
	}
	if patternMatch {
		c.patternHit++
	}
	if !e.Timestamp.IsZero() {
		if !c.earliestSet || e.Timestamp.Before(c.earliestTS) {
			c.earliestTS = e.Timestamp
			c.earliestSet = true
		}
		if e.Timestamp.After(c.latestTS) {
			c.latestTS = e.Timestamp
		}
	}
}

// Summary returns a snapshot of the collected statistics.
func (c *Collector) Summary() Summary {
	return Summary{
		Total:      c.total,
		Matched:    c.matched,
		Skipped:    c.skipped,
		Earliest:   c.earliestTS,
		Latest:     c.latestTS,
		PatternHit: c.patternHit,
	}
}
