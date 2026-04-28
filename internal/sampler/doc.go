// Package sampler provides entry-stream sampling for logslice.
//
// Two strategies are supported:
//
//   - ModeNth  – keeps every Nth log entry, useful for thinning high-volume
//     streams while preserving an even distribution.
//
//   - ModeRandom – keeps each entry independently with a given probability,
//     producing a statistically uniform random sample.
//
// Usage:
//
//	s, _ := sampler.New(sampler.ModeNth, 10, 0, nil)
//	s.Sample(entryCh, outCh)
//
// The zero-based counter used by ModeNth resets after every kept entry so
// the sampling ratio stays constant across arbitrarily long streams.
package sampler
