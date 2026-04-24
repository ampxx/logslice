package output

import (
	"fmt"
	"strings"
)

// FormatFlag is a flag.Value-compatible wrapper around Format,
// suitable for use with the standard flag or pflag packages.
type FormatFlag struct {
	Value Format
}

// String returns the current format value as a string.
func (f *FormatFlag) String() string {
	if f.Value == "" {
		return string(FormatDefault)
	}
	return string(f.Value)
}

// Set parses and validates a format string provided via CLI flag.
func (f *FormatFlag) Set(s string) error {
	switch Format(strings.ToLower(s)) {
	case FormatDefault, FormatJSON, FormatTimestamp:
		f.Value = Format(strings.ToLower(s))
		return nil
	default:
		return fmt.Errorf("unknown format %q: must be one of default, json, timestamp", s)
	}
}

// Type returns the type name used in flag usage messages.
func (f *FormatFlag) Type() string {
	return "format"
}

// Formats returns all supported Format values.
func Formats() []Format {
	return []Format{FormatDefault, FormatJSON, FormatTimestamp}
}
