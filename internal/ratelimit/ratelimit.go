// Package ratelimit provides a token-bucket style rate limiter that caps
// the number of log entries emitted per second. This is useful when tailing
// high-volume log streams and only a representative sample is needed.
package ratelimit

import (
	"errors"
	"time"

	"github.com/mitchellh/logslice/internal/parser"
)

// ErrInvalidRate is returned when the requested rate is not positive.
var ErrInvalidRate = errors.New("ratelimit: rate must be greater than zero")

// Limiter drops log entries that exceed a maximum number of entries per second.
type Limiter struct {
	rate     int
	bucket   int
	last     time.Time
	nowFunc  func() time.Time
}

// New creates a Limiter that allows at most rate entries per second.
// rate must be >= 1.
func New(rate int) (*Limiter, error) {
	if rate <= 0 {
		return nil, ErrInvalidRate
	}
	return &Limiter{
		rate:    rate,
		bucket:  rate,
		last:    time.Time{},
		nowFunc: time.Now,
	}, nil
}

// Allow reports whether the entry should be forwarded.
// It refills the token bucket based on elapsed time since the last call.
func (l *Limiter) Allow(_ parser.Entry) bool {
	now := l.nowFunc()

	if !l.last.IsZero() {
		elapsed := now.Sub(l.last).Seconds()
		refill := int(elapsed * float64(l.rate))
		if refill > 0 {
			l.bucket += refill
			if l.bucket > l.rate {
				l.bucket = l.rate
			}
			l.last = now
		}
	} else {
		l.last = now
	}

	if l.bucket > 0 {
		l.bucket--
		return true
	}
	return false
}

// Apply filters entries from in, forwarding only those permitted by the limiter.
func (l *Limiter) Apply(in []parser.Entry) []parser.Entry {
	out := make([]parser.Entry, 0, len(in))
	for _, e := range in {
		if l.Allow(e) {
			out = append(out, e)
		}
	}
	return out
}
