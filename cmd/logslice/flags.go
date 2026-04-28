package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/nickpoorman/logslice/internal/output"
	"github.com/nickpoorman/logslice/internal/sampler"
)

// config holds all parsed CLI options.
type config struct {
	files      []string
	pattern    string
	from       time.Time
	to         time.Time
	format     output.Format
	stats      bool
	tail       bool
	noColor    bool
	// sampling
	sampleNth  int
	sampleProb float64
}

func parseFlags(args []string) (config, error) {
	fs := flag.NewFlagSet("logslice", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	var (
		pattern    = fs.String("pattern", "", "regex pattern to filter log lines")
		fromStr    = fs.String("from", "", "start time (RFC3339)")
		toStr      = fs.String("to", "", "end time (RFC3339)")
		formatStr  = fs.String("format", "default", fmt.Sprintf("output format %v", output.Formats))
		statsFlag  = fs.Bool("stats", false, "print summary statistics after output")
		tailFlag   = fs.Bool("tail", false, "tail file for new lines")
		noColor    = fs.Bool("no-color", false, "disable ANSI color highlighting")
		sampleNth  = fs.Int("sample-nth", 0, "keep every Nth entry (0 = disabled)")
		sampleProb = fs.Float64("sample-prob", 0, "keep each entry with probability p in (0,1] (0 = disabled)")
	)

	if err := fs.Parse(args); err != nil {
		return config{}, err
	}

	cfg := config{
		files:      fs.Args(),
		pattern:    *pattern,
		stats:      *statsFlag,
		tail:       *tailFlag,
		noColor:    *noColor,
		sampleNth:  *sampleNth,
		sampleProb: *sampleProb,
	}

	if *fromStr != "" {
		t, err := time.Parse(time.RFC3339, *fromStr)
		if err != nil {
			return config{}, fmt.Errorf("invalid --from: %w", err)
		}
		cfg.from = t
	}
	if *toStr != "" {
		t, err := time.Parse(time.RFC3339, *toStr)
		if err != nil {
			return config{}, fmt.Errorf("invalid --to: %w", err)
		}
		cfg.to = t
	}

	fmt, err := output.ParseFormat(*formatStr)
	if err != nil {
		return config{}, err
	}
	cfg.format = fmt

	if *sampleNth > 0 && *sampleProb > 0 {
		return config{}, errors.New("--sample-nth and --sample-prob are mutually exclusive")
	}
	if *sampleNth > 0 {
		if _, err := sampler.New(sampler.ModeNth, *sampleNth, 0, nil); err != nil {
			return config{}, fmt.Errorf("invalid --sample-nth: %w", err)
		}
	}
	if *sampleProb > 0 {
		if _, err := sampler.New(sampler.ModeRandom, 0, *sampleProb, nil); err != nil {
			return config{}, fmt.Errorf("invalid --sample-prob: %w", err)
		}
	}

	return cfg, nil
}
