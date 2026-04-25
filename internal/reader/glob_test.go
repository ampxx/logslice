package reader

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func writeNamedTempFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("writeNamedTempFile: %v", err)
	}
	return path
}

func TestGlobSource_ReadsMultipleFiles(t *testing.T) {
	dir := t.TempDir()
	writeNamedTempFile(t, dir, "a.log", "line-a1\nline-a2\n")
	writeNamedTempFile(t, dir, "b.log", "line-b1\n")

	pattern := filepath.Join(dir, "*.log")
	ch, name, err := GlobSource(pattern)
	if err != nil {
		t.Fatalf("GlobSource error: %v", err)
	}
	if name == "" {
		t.Error("expected non-empty source name")
	}

	var lines []string
	for l := range ch {
		lines = append(lines, l)
	}

	want := []string{"line-a1", "line-a2", "line-b1"}
	if len(lines) != len(want) {
		t.Fatalf("got %d lines, want %d: %v", len(lines), len(want), lines)
	}
	for i, w := range want {
		if lines[i] != w {
			t.Errorf("line[%d]: got %q, want %q", i, lines[i], w)
		}
	}
}

func TestGlobSource_NoMatch(t *testing.T) {
	dir := t.TempDir()
	_, _, err := GlobSource(filepath.Join(dir, "*.log"))
	if err == nil {
		t.Fatal("expected error for non-matching glob, got nil")
	}
}

func TestGlobSource_InvalidPattern(t *testing.T) {
	_, _, err := GlobSource("[invalid")
	if err == nil {
		t.Fatal("expected error for invalid glob pattern, got nil")
	}
}

func TestMultiFileSource_OrderPreserved(t *testing.T) {
	dir := t.TempDir()
	paths := make([]string, 3)
	for i := range paths {
		paths[i] = writeNamedTempFile(t, dir, fmt.Sprintf("f%d.log", i), fmt.Sprintf("line%d\n", i))
	}

	ch, name, err := MultiFileSource(paths)
	if err != nil {
		t.Fatalf("MultiFileSource error: %v", err)
	}
	if name == "" {
		t.Error("expected non-empty source name")
	}

	var lines []string
	for l := range ch {
		lines = append(lines, l)
	}

	if len(lines) != 3 {
		t.Fatalf("got %d lines, want 3", len(lines))
	}
	for i, want := range []string{"line0", "line1", "line2"} {
		if lines[i] != want {
			t.Errorf("line[%d]: got %q, want %q", i, lines[i], want)
		}
	}
}

func TestMultiFileSource_NoPaths(t *testing.T) {
	_, _, err := MultiFileSource(nil)
	if err == nil {
		t.Fatal("expected error for empty paths, got nil")
	}
}
