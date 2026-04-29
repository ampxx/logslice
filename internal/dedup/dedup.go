// Package dedup provides log entry deduplication by detecting and suppressing
// consecutive or near-consecutive repeated log lines.
package dedup

import (
	"hash/fnv"
	"strconv"

	"github.com/example/logslice/internal/parser"
)

// Mode controls how deduplication is applied.
type Mode int

const (
	// Consecutive suppresses only back-to-back duplicate lines.
	Consecutive Mode = iota
	// Window suppresses duplicates seen within a sliding window of N entries.
	Window
)

// Deduplicator filters repeated log entries.
type Deduplicator struct {
	mode    Mode
	window  int
	seen    []uint64
	pos     int
	count   int
	Skipped int
}

// New creates a Deduplicator. windowSize is only used in Window mode;
// pass 1 for Consecutive mode.
func New(mode Mode, windowSize int) (*Deduplicator, error) {
	if windowSize < 1 {
		return nil, ErrInvalidWindow
	}
	return &Deduplicator{
		mode:   mode,
		window: windowSize,
		seen:   make([]uint64, windowSize),
	}, nil
}

// IsDuplicate returns true if the entry is a duplicate according to the
// configured mode, and increments the internal skip counter.
func (d *Deduplicator) IsDuplicate(e parser.Entry) bool {
	h := hashEntry(e)
	if d.mode == Consecutive {
		if d.count > 0 && d.seen[0] == h {
			d.Skipped++
			return true
		}
		d.seen[0] = h
		d.count++
		return false
	}
	// Window mode
	for i := 0; i < d.count && i < d.window; i++ {
		if d.seen[(d.pos-1-i+d.window)%d.window] == h {
			d.Skipped++
			return true
		}
	}
	d.seen[d.pos] = h
	d.pos = (d.pos + 1) % d.window
	if d.count < d.window {
		d.count++
	}
	return false
}

func hashEntry(e parser.Entry) uint64 {
	h := fnv.New64a()
	_, _ = h.Write([]byte(e.Raw))
	if !e.Timestamp.IsZero() {
		_, _ = h.Write([]byte(strconv.FormatInt(e.Timestamp.UnixNano(), 10)))
	}
	return h.Sum64()
}
