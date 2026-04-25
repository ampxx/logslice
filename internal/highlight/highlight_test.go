package highlight_test

import (
	"strings"
	"testing"

	"github.com/yourorg/logslice/internal/highlight"
)

func TestHighlighter_NoPattern(t *testing.T) {
	h, err := highlight.New("", highlight.Yellow)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	line := "nothing should change"
	if got := h.Apply(line); got != line {
		t.Errorf("expected %q, got %q", line, got)
	}
}

func TestHighlighter_InvalidPattern(t *testing.T) {
	_, err := highlight.New("[", highlight.Red)
	if err == nil {
		t.Fatal("expected error for invalid pattern, got nil")
	}
}

func TestHighlighter_SingleMatch(t *testing.T) {
	h, err := highlight.New("ERROR", highlight.Red)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	line := "2024-01-01 ERROR something failed"
	got := h.Apply(line)
	if !strings.Contains(got, highlight.Red+"ERROR"+highlight.Reset) {
		t.Errorf("expected highlighted ERROR in %q", got)
	}
	if !strings.Contains(got, "something failed") {
		t.Errorf("expected remainder of line preserved in %q", got)
	}
}

func TestHighlighter_MultipleMatches(t *testing.T) {
	h, err := highlight.New("foo", highlight.Cyan)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	line := "foo bar foo baz foo"
	got := h.Apply(line)
	count := strings.Count(got, highlight.Cyan+"foo"+highlight.Reset)
	if count != 3 {
		t.Errorf("expected 3 highlighted matches, got %d in %q", count, got)
	}
}

func TestHighlighter_NoMatch(t *testing.T) {
	h, err := highlight.New("WARN", highlight.Yellow)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	line := "INFO everything is fine"
	got := h.Apply(line)
	if got != line {
		t.Errorf("expected unchanged line %q, got %q", line, got)
	}
}

func TestColor(t *testing.T) {
	got := highlight.Color(highlight.Bold, "hello")
	if !strings.HasPrefix(got, highlight.Bold) {
		t.Errorf("expected bold prefix in %q", got)
	}
	if !strings.HasSuffix(got, highlight.Reset) {
		t.Errorf("expected reset suffix in %q", got)
	}
}
