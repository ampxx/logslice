package ratelimit

import (
	"testing"
	"time"

	"github.com/mitchellh/logslice/internal/parser"
)

func makeEntry(raw string) parser.Entry {
	return parser.Entry{Raw: raw}
}

func TestNew_InvalidRate(t *testing.T) {
	_, err := New(0)
	if err == nil {
		t.Fatal("expected error for rate=0")
	}
	_, err = New(-5)
	if err == nil {
		t.Fatal("expected error for negative rate")
	}
}

func TestNew_ValidRate(t *testing.T) {
	l, err := New(10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if l.rate != 10 {
		t.Fatalf("expected rate 10, got %d", l.rate)
	}
}

func TestAllow_BucketDrains(t *testing.T) {
	l, _ := New(3)
	fixedNow := time.Now()
	l.nowFunc = func() time.Time { return fixedNow }

	allowed := 0
	for i := 0; i < 6; i++ {
		if l.Allow(makeEntry("line")) {
			allowed++
		}
	}
	if allowed != 3 {
		t.Fatalf("expected 3 allowed (bucket size), got %d", allowed)
	}
}

func TestAllow_BucketRefillsOverTime(t *testing.T) {
	l, _ := New(5)
	base := time.Now()
	call := 0
	times := []time.Time{
		base,
		base,
		base,
		base,
		base,
		base,                        // bucket drained after 5
		base.Add(2 * time.Second),   // +2 s => refill 10, capped at 5
		base.Add(2 * time.Second),
	}
	l.nowFunc = func() time.Time {
		t := times[call]
		if call < len(times)-1 {
			call++
		}
		return t
	}

	// drain the bucket
	for i := 0; i < 6; i++ {
		l.Allow(makeEntry("x"))
	}
	// after refill both should pass
	if !l.Allow(makeEntry("x")) {
		t.Fatal("expected entry allowed after refill")
	}
	if !l.Allow(makeEntry("x")) {
		t.Fatal("expected second entry allowed after refill")
	}
}

func TestApply_FiltersSlice(t *testing.T) {
	l, _ := New(2)
	fixedNow := time.Now()
	l.nowFunc = func() time.Time { return fixedNow }

	entries := []parser.Entry{
		makeEntry("a"), makeEntry("b"), makeEntry("c"), makeEntry("d"),
	}
	out := l.Apply(entries)
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
	if out[0].Raw != "a" || out[1].Raw != "b" {
		t.Fatalf("unexpected entries: %+v", out)
	}
}
