package reader_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/logslice/internal/reader"
)

func writeTempFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "test.log")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("writeTempFile: %v", err)
	}
	return path
}

func TestNewFileSource_ReadsLines(t *testing.T) {
	path := writeTempFile(t, "line one\nline two\nline three\n")

	src, err := reader.NewFileSource(path)
	if err != nil {
		t.Fatalf("NewFileSource: %v", err)
	}
	defer src.Close()

	var got []string
	for line := range src.Lines() {
		got = append(got, line)
	}

	want := []string{"line one", "line two", "line three"}
	if len(got) != len(want) {
		t.Fatalf("got %d lines, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("line %d: got %q, want %q", i, got[i], want[i])
		}
	}
}

func TestNewFileSource_EmptyFile(t *testing.T) {
	path := writeTempFile(t, "")

	src, err := reader.NewFileSource(path)
	if err != nil {
		t.Fatalf("NewFileSource: %v", err)
	}
	defer src.Close()

	var count int
	for range src.Lines() {
		count++
	}
	if count != 0 {
		t.Errorf("expected 0 lines, got %d", count)
	}
}

func TestNewFileSource_MissingFile(t *testing.T) {
	_, err := reader.NewFileSource("/nonexistent/path/to/file.log")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestNewStdinSource_Name(t *testing.T) {
	src := reader.NewStdinSource()
	if src.Name != "<stdin>" {
		t.Errorf("Name = %q, want %q", src.Name, "<stdin>")
	}
}
