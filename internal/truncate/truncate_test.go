package truncate_test

import (
	"testing"
	"time"

	"github.com/dkoosis/logslice/internal/parser"
	"github.com/dkoosis/logslice/internal/truncate"
)

func makeEntry(raw string) parser.Entry {
	return parser.Entry{
		Timestamp: time.Time{},
		Raw:       raw,
	}
}

func TestNew_InvalidMaxLen(t *testing.T) {
	_, err := truncate.New(0, "…")
	if err == nil {
		t.Fatal("expected error for maxLen=0, got nil")
	}
	_, err = truncate.New(-5, "…")
	if err == nil {
		t.Fatal("expected error for maxLen=-5, got nil")
	}
}

func TestNew_ValidMaxLen(t *testing.T) {
	tr, err := truncate.New(10, "…")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tr == nil {
		t.Fatal("expected non-nil Truncator")
	}
}

func TestApply_ShortLine_Unchanged(t *testing.T) {
	tr, _ := truncate.New(20, "…")
	e := makeEntry("hello world")
	got := tr.Apply(e)
	if got.Raw != "hello world" {
		t.Errorf("expected unchanged raw, got %q", got.Raw)
	}
}

func TestApply_ExactLength_Unchanged(t *testing.T) {
	tr, _ := truncate.New(5, "…")
	e := makeEntry("abcde")
	got := tr.Apply(e)
	if got.Raw != "abcde" {
		t.Errorf("expected %q, got %q", "abcde", got.Raw)
	}
}

func TestApply_LongLine_Truncated(t *testing.T) {
	tr, _ := truncate.New(10, "…")
	e := makeEntry("this is a very long log line that exceeds the limit")
	got := tr.Apply(e)
	// suffix is 1 rune, so cut = 9 runes of text + "…"
	expected := "this is a…"
	if got.Raw != expected {
		t.Errorf("expected %q, got %q", expected, got.Raw)
	}
}

func TestApply_NoSuffix(t *testing.T) {
	tr, _ := truncate.New(5, "")
	e := makeEntry("abcdefgh")
	got := tr.Apply(e)
	if got.Raw != "abcde" {
		t.Errorf("expected %q, got %q", "abcde", got.Raw)
	}
}

func TestApply_PreservesTimestamp(t *testing.T) {
	tr, _ := truncate.New(5, "…")
	now := time.Now()
	e := parser.Entry{Timestamp: now, Raw: "longer than five chars"}
	got := tr.Apply(e)
	if !got.Timestamp.Equal(now) {
		t.Errorf("timestamp changed: got %v, want %v", got.Timestamp, now)
	}
}

func TestApplyAll_TruncatesAll(t *testing.T) {
	tr, _ := truncate.New(4, "-")
	entries := []parser.Entry{
		makeEntry("short"),
		makeEntry("toolong"),
		makeEntry("ok"),
	}
	out := tr.ApplyAll(entries)
	if len(out) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(out))
	}
	if out[1].Raw != "too-" {
		t.Errorf("expected %q, got %q", "too-", out[1].Raw)
	}
	if out[2].Raw != "ok" {
		t.Errorf("expected unchanged %q, got %q", "ok", out[2].Raw)
	}
}
