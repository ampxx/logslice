package levelfilter_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/levelfilter"
	"github.com/yourorg/logslice/internal/parser"
)

func makeEntry(raw string) parser.Entry {
	return parser.Entry{Raw: raw}
}

func TestParseLevel_KnownTokens(t *testing.T) {
	cases := []struct {
		input string
		want  levelfilter.Level
	}{
		{"debug", levelfilter.LevelDebug},
		{"INFO", levelfilter.LevelInfo},
		{"WARN", levelfilter.LevelWarn},
		{"warning", levelfilter.LevelWarn},
		{"ERROR", levelfilter.LevelError},
		{"err", levelfilter.LevelError},
		{"FATAL", levelfilter.LevelFatal},
		{"critical", levelfilter.LevelFatal},
	}
	for _, tc := range cases {
		got, ok := levelfilter.ParseLevel(tc.input)
		if !ok {
			t.Errorf("ParseLevel(%q): expected ok=true", tc.input)
		}
		if got != tc.want {
			t.Errorf("ParseLevel(%q) = %d, want %d", tc.input, got, tc.want)
		}
	}
}

func TestParseLevel_Unknown(t *testing.T) {
	_, ok := levelfilter.ParseLevel("verbose")
	if ok {
		t.Error("expected ok=false for unknown level")
	}
}

func TestFilter_MinLevelWarn(t *testing.T) {
	entries := []parser.Entry{
		makeEntry("2024-01-01 DEBUG starting up"),
		makeEntry("2024-01-01 INFO  server ready"),
		makeEntry("2024-01-01 WARN  disk space low"),
		makeEntry("2024-01-01 ERROR connection refused"),
		makeEntry("2024-01-01 FATAL out of memory"),
	}
	got := levelfilter.Filter(entries, levelfilter.LevelWarn)
	if len(got) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(got))
	}
}

func TestFilter_NoLevelEntryAlwaysKept(t *testing.T) {
	entries := []parser.Entry{
		makeEntry("plain log line with no level token"),
		makeEntry("another unstructured line"),
	}
	got := levelfilter.Filter(entries, levelfilter.LevelError)
	if len(got) != 2 {
		t.Fatalf("expected 2 entries (no level = keep), got %d", len(got))
	}
}

func TestFilter_AllLevelsPassDebug(t *testing.T) {
	entries := []parser.Entry{
		makeEntry("DEBUG trace"),
		makeEntry("INFO  msg"),
		makeEntry("ERROR boom"),
	}
	got := levelfilter.Filter(entries, levelfilter.LevelDebug)
	if len(got) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(got))
	}
}
