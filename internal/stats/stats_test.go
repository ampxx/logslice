package stats

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/parser"
)

func makeEntry(ts time.Time, raw string) parser.Entry {
	return parser.Entry{Timestamp: ts, Raw: raw}
}

func TestCollector_BasicCounts(t *testing.T) {
	c := NewCollector()
	now := time.Now()
	c.Record(makeEntry(now, "line1"), true, true)
	c.Record(makeEntry(now.Add(time.Minute), "line2"), true, false)
	c.Record(makeEntry(time.Time{}, "line3"), false, false)

	s := c.Summary()
	if s.Total != 3 {
		t.Errorf("Total: want 3, got %d", s.Total)
	}
	if s.Matched != 2 {
		t.Errorf("Matched: want 2, got %d", s.Matched)
	}
	if s.Skipped != 1 {
		t.Errorf("Skipped: want 1, got %d", s.Skipped)
	}
	if s.PatternHit != 1 {
		t.Errorf("PatternHit: want 1, got %d", s.PatternHit)
	}
}

func TestCollector_TimestampRange(t *testing.T) {
	c := NewCollector()
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	c.Record(makeEntry(base.Add(2*time.Hour), "b"), true, false)
	c.Record(makeEntry(base, "a"), true, false)
	c.Record(makeEntry(base.Add(5*time.Hour), "c"), true, false)

	s := c.Summary()
	if !s.Earliest.Equal(base) {
		t.Errorf("Earliest: want %v, got %v", base, s.Earliest)
	}
	if !s.Latest.Equal(base.Add(5 * time.Hour)) {
		t.Errorf("Latest: want %v, got %v", base.Add(5*time.Hour), s.Latest)
	}
}

func TestCollector_NoTimestamps(t *testing.T) {
	c := NewCollector()
	c.Record(makeEntry(time.Time{}, "x"), true, false)
	s := c.Summary()
	if !s.Earliest.IsZero() {
		t.Error("expected zero Earliest when no timestamps")
	}
}

func TestPrint_ContainsExpectedFields(t *testing.T) {
	base := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	s := Summary{
		Total: 10, Matched: 8, Skipped: 2,
		PatternHit: 3,
		Earliest:   base,
		Latest:     base.Add(30 * time.Minute),
	}
	var buf bytes.Buffer
	Print(&buf, s)
	out := buf.String()
	for _, want := range []string{"10", "8", "2", "3", "30m0s"} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q:\n%s", want, out)
		}
	}
}
