package filter_test

import (
	"regexp"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/filter"
	"github.com/yourorg/logslice/internal/parser"
)

func makeEntry(ts time.Time, raw string) parser.Entry {
	return parser.Entry{Timestamp: ts, Raw: raw}
}

var (
	t0 = time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	t1 = time.Date(2024, 1, 1, 11, 0, 0, 0, time.UTC)
	t2 = time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
)

func TestFilter_TimeRange(t *testing.T) {
	entry := makeEntry(t1, "2024-01-01 11:00:00 INFO hello")

	if !filter.Filter(entry, filter.Options{From: t0, To: t2}) {
		t.Error("expected entry within range to pass")
	}
	if filter.Filter(entry, filter.Options{From: t2}) {
		t.Error("expected entry before From to be rejected")
	}
	if filter.Filter(entry, filter.Options{To: t0}) {
		t.Error("expected entry after To to be rejected")
	}
}

func TestFilter_Pattern(t *testing.T) {
	entry := makeEntry(t1, "2024-01-01 11:00:00 ERROR disk full")

	pat := regexp.MustCompile(`ERROR`)
	if !filter.Filter(entry, filter.Options{Pattern: pat}) {
		t.Error("expected pattern match to pass")
	}

	noMatch := regexp.MustCompile(`WARN`)
	if filter.Filter(entry, filter.Options{Pattern: noMatch}) {
		t.Error("expected non-matching pattern to be rejected")
	}
}

func TestFilter_NoOptions(t *testing.T) {
	entry := makeEntry(t1, "some log line")
	if !filter.Filter(entry, filter.Options{}) {
		t.Error("expected entry to pass with no filter options")
	}
}

func TestFilterAll(t *testing.T) {
	entries := []parser.Entry{
		makeEntry(t0, "line at t0"),
		makeEntry(t1, "line at t1"),
		makeEntry(t2, "line at t2"),
	}

	result := filter.FilterAll(entries, filter.Options{From: t1, To: t2})
	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
	if !result[0].Timestamp.Equal(t1) {
		t.Errorf("unexpected first entry timestamp: %v", result[0].Timestamp)
	}
}
