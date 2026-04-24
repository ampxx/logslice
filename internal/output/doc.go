// Package output provides formatting and rendering utilities for parsed log
// entries produced by the parser and filter packages.
//
// # Formats
//
// Three output formats are supported:
//
//   - FormatDefault  – emits the original raw log line unchanged.
//   - FormatJSON     – emits a compact JSON object with "timestamp" and "raw" fields.
//   - FormatTimestamp – prefixes each line with an RFC3339 timestamp.
//
// # Usage
//
//	f := output.NewFormatter(os.Stdout, output.FormatJSON)
//	f.WriteAll(filteredEntries)
package output
