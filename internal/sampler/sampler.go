package sampler

import (
	"math/rand"

	"github.com/nickpoorman/logslice/internal/parser"
)

// Mode controls how sampling is applied.
type Mode int

const (
	// ModeNth emits every Nth entry.
	ModeNth Mode = iota
	// ModeRandom emits each entry with a given probability.
	ModeRandom
)

// Sampler reduces a stream of log entries according to a sampling strategy.
type Sampler struct {
	mode Mode
	n    int     // used by ModeNth
	prob float64 // used by ModeRandom, in range (0, 1]
	counter int
	rng     *rand.Rand
}

// New creates a Sampler. For ModeNth supply n >= 1; for ModeRandom supply
// probability in (0, 1].
func New(mode Mode, n int, prob float64, src rand.Source) (*Sampler, error) {
	if mode == ModeNth && n < 1 {
		return nil, ErrInvalidN
	}
	if mode == ModeRandom && (prob <= 0 || prob > 1) {
		return nil, ErrInvalidProb
	}
	if src == nil {
		src = rand.NewSource(42)
	}
	return &Sampler{
		mode: mode,
		n:    n,
		prob: prob,
		rng:  rand.New(src),
	}, nil
}

// Keep reports whether the entry should be kept in the output stream.
func (s *Sampler) Keep(_ parser.Entry) bool {
	switch s.mode {
	case ModeNth:
		s.counter++
		if s.counter >= s.n {
			s.counter = 0
			return true
		}
		return false
	case ModeRandom:
		return s.rng.Float64() < s.prob
	}
	return true
}

// Sample filters entries from in, writing kept entries to out.
func (s *Sampler) Sample(in <-chan parser.Entry, out chan<- parser.Entry) {
	for e := range in {
		if s.Keep(e) {
			out <- e
		}
	}
}
