// Package stats provides lightweight aggregation of log entry statistics
// collected during a logslice run.
//
// Usage:
//
//	collector := stats.NewCollector()
//
//	for _, entry := range entries {
//		matched := filter.Filter(entry, opts)
//		collector.Record(entry, matched, patternMatched)
//	}
//
//	summary := collector.Summary()
//	stats.Print(os.Stderr, summary)
//
// The Collector is not safe for concurrent use; callers that process entries
// in parallel should use separate collectors and merge results manually.
package stats
