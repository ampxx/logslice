// Package levelfilter implements severity-level based filtering for log entries.
//
// It recognises common level tokens (DEBUG, INFO, WARN, WARNING, ERROR, FATAL,
// and their common abbreviations) embedded anywhere in a raw log line and
// exposes a Filter function that discards entries below a caller-specified
// minimum severity.
//
// Usage:
//
//	minLvl, ok := levelfilter.ParseLevel(flagValue)
//	if !ok {
//		log.Fatalf("unknown log level: %s", flagValue)
//	}
//	filtered := levelfilter.Filter(entries, minLvl)
//
// Entries whose raw text contains no recognisable level token are treated as
// opaque and are always passed through, preserving lines that use a custom or
// non-standard level vocabulary.
package levelfilter
