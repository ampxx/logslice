package dedup

import (
	"testing"
	"time"

	"github.com/example/logslice/internal/parser"
)

func makeEntry(raw string) parser.Entry {
	return parser.Entry{Raw: raw, Timestamp: time.Time{}}
}

func TestNew_InvalidWindow(t *testing.T) {
	_, err := New(Window, 0)
	if err != ErrInvalidWindow {
		t.Fatalf("expected ErrInvalidWindow, got %v", err)
	}
}

func TestNew_ValidWindow(t *testing.T) {
	d, err := New(Window, 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d == nil {
		t.Fatal("expected non-nil deduplicator")
	}
}

func TestConsecutive_NoDupe(t *testing.T) {
	d, _ := New(Consecutive, 1)
	if d.IsDuplicate(makeEntry("line one")) {
		t.Error("first entry should not be duplicate")
	}
	if d.IsDuplicate(makeEntry("line two")) {
		t.Error("different entry should not be duplicate")
	}
}

func TestConsecutive_DetectsDupe(t *testing.T) {
	d, _ := New(Consecutive, 1)
	d.IsDuplicate(makeEntry("repeated line"))
	if !d.IsDuplicate(makeEntry("repeated line")) {
		t.Error("identical consecutive entry should be duplicate")
	}
	if d.Skipped != 1 {
		t.Errorf("expected 1 skipped, got %d", d.Skipped)
	}
}

func TestConsecutive_ResetOnDifferent(t *testing.T) {
	d, _ := New(Consecutive, 1)
	d.IsDuplicate(makeEntry("aaa"))
	d.IsDuplicate(makeEntry("bbb"))
	if d.IsDuplicate(makeEntry("aaa")) {
		t.Error("non-consecutive repeat should not be duplicate in Consecutive mode")
	}
}

func TestWindow_DetectsDupeInWindow(t *testing.T) {
	d, _ := New(Window, 4)
	entries := []string{"a", "b", "c", "a"}
	results := make([]bool, len(entries))
	for i, raw := range entries {
		results[i] = d.IsDuplicate(makeEntry(raw))
	}
	if results[3] != true {
		t.Error("'a' should be duplicate within window of 4")
	}
	if d.Skipped != 1 {
		t.Errorf("expected 1 skipped, got %d", d.Skipped)
	}
}

func TestWindow_NoDupeOutsideWindow(t *testing.T) {
	d, _ := New(Window, 2)
	d.IsDuplicate(makeEntry("x"))
	d.IsDuplicate(makeEntry("y"))
	d.IsDuplicate(makeEntry("z"))
	if d.IsDuplicate(makeEntry("x")) {
		t.Error("'x' should not be duplicate — it fell outside the window")
	}
}
