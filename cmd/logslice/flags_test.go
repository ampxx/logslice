package main

import (
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/output"
)

func TestParseFlags_Defaults(t *testing.T) {
	cfg, err := parseFlags([]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.inputFile != "" {
		t.Errorf("expected empty inputFile, got %q", cfg.inputFile)
	}
	if cfg.pattern != "" {
		t.Errorf("expected empty pattern, got %q", cfg.pattern)
	}
	if cfg.format != output.FormatDefault {
		t.Errorf("expected default format, got %v", cfg.format)
	}
}

func TestParseFlags_AllOptions(t *testing.T) {
	args := []string{
		"-f", "app.log",
		"-p", "ERROR",
		"-o", "json",
		"-from", "2024-01-01T00:00:00Z",
		"-to", "2024-01-02T00:00:00Z",
	}
	cfg, err := parseFlags(args)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.inputFile != "app.log" {
		t.Errorf("expected app.log, got %q", cfg.inputFile)
	}
	if cfg.pattern != "ERROR" {
		t.Errorf("expected ERROR pattern, got %q", cfg.pattern)
	}
	if cfg.format != output.FormatJSON {
		t.Errorf("expected json format, got %v", cfg.format)
	}
	wantFrom, _ := time.Parse(time.RFC3339, "2024-01-01T00:00:00Z")
	if !cfg.from.Equal(wantFrom) {
		t.Errorf("unexpected from time: %v", cfg.from)
	}
}

func TestParseFlags_InvalidFrom(t *testing.T) {
	_, err := parseFlags([]string{"-from", "not-a-time"})
	if err == nil {
		t.Error("expected error for invalid -from value")
	}
}

func TestParseFlags_InvalidFormat(t *testing.T) {
	_, err := parseFlags([]string{"-o", "xml"})
	if err == nil {
		t.Error("expected error for unknown format")
	}
}
