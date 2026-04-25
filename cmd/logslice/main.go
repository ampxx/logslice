package main

import (
	"fmt"
	"os"

	"github.com/yourorg/logslice/internal/filter"
	"github.com/yourorg/logslice/internal/output"
	"github.com/yourorg/logslice/internal/parser"
	"github.com/yourorg/logslice/internal/reader"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "logslice: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := parseFlags(os.Args[1:])
	if err != nil {
		return err
	}

	var src reader.Source
	if cfg.inputFile == "" || cfg.inputFile == "-" {
		src = reader.NewStdinSource()
	} else {
		src, err = reader.NewFileSource(cfg.inputFile)
		if err != nil {
			return err
		}
	}

	lp := parser.NewLineParser()
	fmt, err := output.NewFormatter(cfg.format)
	if err != nil {
		return err
	}

	opts := filter.Options{
		From:    cfg.from,
		To:      cfg.to,
		Pattern: cfg.pattern,
	}

	for line := range src.Lines() {
		entry, ok := lp.Parse(line)
		if !ok {
			continue
		}
		if !filter.Filter(entry, opts) {
			continue
		}
		if err := fmt.Write(os.Stdout, entry); err != nil {
			return err
		}
	}
	return src.Err()
}
