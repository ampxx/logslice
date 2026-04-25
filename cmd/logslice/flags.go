package main

import (
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/yourorg/logslice/internal/output"
)

type config struct {
	inputFile string
	pattern   string
	format    output.Format
	from      time.Time
	to        time.Time
}

func parseFlags(args []string) (*config, error) {
	fs := flag.NewFlagSet("logslice", flag.ContinueOnError)

	var (
		inputFile  = fs.String("f", "", "input log file (default: stdin)")
		pattern    = fs.String("p", "", "regex pattern to match log lines")
		formatStr  = fs.String("o", "default", fmt.Sprintf("output format: %s", strings.Join(output.Formats(), ", ")))
		fromStr    = fs.String("from", "", "start time filter (RFC3339, e.g. 2024-01-01T00:00:00Z)")
		toStr      = fs.String("to", "", "end time filter (RFC3339)")
	)

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	fmt, err := output.ParseFormat(*formatStr)
	if err != nil {
		return nil, err
	}

	cfg := &config{
		inputFile: *inputFile,
		pattern:   *pattern,
		format:    fmt,
	}

	if *fromStr != "" {
		t, err := time.Parse(time.RFC3339, *fromStr)
		if err != nil {
			return nil, fmt.Errorf("invalid -from value: %w", err)
		}
		cfg.from = t
	}

	if *toStr != "" {
		t, err := time.Parse(time.RFC3339, *toStr)
		if err != nil {
			return nil, fmt.Errorf("invalid -to value: %w", err)
		}
		cfg.to = t
	}

	return cfg, nil
}
