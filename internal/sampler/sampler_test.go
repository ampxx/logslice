package sampler_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/nickpoorman/logslice/internal/parser"
	"github.com/nickpoorman/logslice/internal/sampler"
)

func makeEntries(n int) []parser.Entry {
	entries := make([]parser.Entry, n)
	for i := range entries {
		entries[i] = parser.Entry{Raw: "line"}
	}
	return entries
}

func TestNew_InvalidN(t *testing.T) {
	_, err := sampler.New(sampler.ModeNth, 0, 0, nil)
	if err == nil {
		t.Fatal("expected error for n=0")
	}
}

func TestNew_InvalidProb(t *testing.T) {
	_, err := sampler.New(sampler.ModeRandom, 0, 1.5, nil)
	if err == nil {
		t.Fatal("expected error for prob > 1")
	}
	_, err = sampler.New(sampler.ModeRandom, 0, 0, nil)
	if err == nil {
		t.Fatal("expected error for prob == 0")
	}
}

func TestSampler_NthMode(t *testing.T) {
	s, err := sampler.New(sampler.ModeNth, 3, 0, nil)
	if err != nil {
		t.Fatal(err)
	}
	entries := makeEntries(9)
	kept := 0
	for _, e := range entries {
		if s.Keep(e) {
			kept++
		}
	}
	if kept != 3 {
		t.Fatalf("expected 3 kept entries, got %d", kept)
	}
}

func TestSampler_RandomMode(t *testing.T) {
	src := rand.NewSource(time.Now().UnixNano())
	s, err := sampler.New(sampler.ModeRandom, 0, 0.5, src)
	if err != nil {
		t.Fatal(err)
	}
	entries := makeEntries(1000)
	kept := 0
	for _, e := range entries {
		if s.Keep(e) {
			kept++
		}
	}
	// With prob=0.5 and 1000 entries, expect roughly 500 ± 100.
	if kept < 350 || kept > 650 {
		t.Fatalf("random sample far from expected: got %d/1000", kept)
	}
}

func TestSampler_SampleChannel(t *testing.T) {
	s, _ := sampler.New(sampler.ModeNth, 2, 0, nil)
	in := make(chan parser.Entry, 10)
	out := make(chan parser.Entry, 10)
	for _, e := range makeEntries(10) {
		in <- e
	}
	close(in)
	s.Sample(in, out)
	close(out)
	count := 0
	for range out {
		count++
	}
	if count != 5 {
		t.Fatalf("expected 5, got %d", count)
	}
}
