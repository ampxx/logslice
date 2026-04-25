// Package main is the entry point for the logslice command-line tool.
//
// logslice reads log lines from a file or stdin, parses timestamps,
// and emits filtered output based on time range and pattern options.
//
// Usage:
//
//	logslice [flags]
//
// Flags:
//
//	-f <file>       Input log file (default: stdin)
//	-p <pattern>    Regex pattern to match against log lines
//	-o <format>     Output format: default, json, timestamp (default: default)
//	-from <time>    Inclusive start time filter (RFC3339)
//	-to   <time>    Inclusive end time filter (RFC3339)
//
// Examples:
//
//	logslice -f app.log -p ERROR -o json
//	logslice -f app.log -from 2024-01-01T00:00:00Z -to 2024-01-01T23:59:59Z
//	cat app.log | logslice -p WARN
package main
