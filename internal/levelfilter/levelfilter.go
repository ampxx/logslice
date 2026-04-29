// Package levelfilter provides log-level based filtering for log entries.
package levelfilter

import (
	"regexp"
	"strings"

	"github.com/yourorg/logslice/internal/parser"
)

// Level represents a log severity level.
type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

var levelNames = map[Level][]string{
	LevelDebug: {"DEBUG", "DBG", "TRACE"},
	LevelInfo:  {"INFO", "INF"},
	LevelWarn:  {"WARN", "WARNING", "WRN"},
	LevelError: {"ERROR", "ERR"},
	LevelFatal: {"FATAL", "CRIT", "CRITICAL"},
}

var levelPattern = regexp.MustCompile(
	`(?i)\b(DEBUG|DBG|TRACE|INFO|INF|WARN|WARNING|WRN|ERROR|ERR|FATAL|CRIT|CRITICAL)\b`,
)

// ParseLevel parses a level string (case-insensitive) into a Level value.
// Returns LevelDebug and false if the string is not recognised.
func ParseLevel(s string) (Level, bool) {
	upper := strings.ToUpper(strings.TrimSpace(s))
	for lvl, names := range levelNames {
		for _, n := range names {
			if n == upper {
				return lvl, true
			}
		}
	}
	return LevelDebug, false
}

// Filter keeps only entries whose detected log level is >= minLevel.
// Entries with no detectable level are always kept.
func Filter(entries []parser.Entry, minLevel Level) []parser.Entry {
	out := make([]parser.Entry, 0, len(entries))
	for _, e := range entries {
		if lvl, ok := detectLevel(e.Raw); !ok || lvl >= minLevel {
			out = append(out, e)
		}
	}
	return out
}

// detectLevel scans the raw log line for the first recognisable level token.
func detectLevel(line string) (Level, bool) {
	match := levelPattern.FindString(line)
	if match == "" {
		return LevelDebug, false
	}
	return ParseLevel(match)
}
