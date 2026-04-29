package redact

import (
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/parser"
)

func makeEntry(msg string) parser.Entry {
	return parser.Entry{Timestamp: time.Now(), Message: msg, Raw: msg}
}

func TestNew_EmptyPlaceholder(t *testing.T) {
	_, err := New([]string{`\d+`}, "")
	if err != ErrEmptyPlaceholder {
		t.Fatalf("expected ErrEmptyPlaceholder, got %v", err)
	}
}

func TestNew_InvalidPattern(t *testing.T) {
	_, err := New([]string{`[invalid`}, "[REDACTED]")
	if err == nil {
		t.Fatal("expected error for invalid regex, got nil")
	}
}

func TestNew_Valid(t *testing.T) {
	r, err := New([]string{`\d+`}, "[REDACTED]")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r == nil {
		t.Fatal("expected non-nil Redactor")
	}
}

func TestApply_NoPatterns(t *testing.T) {
	r, _ := New(nil, "[REDACTED]")
	e := makeEntry("hello world")
	out := r.Apply(e)
	if out.Message != "hello world" {
		t.Errorf("expected unchanged message, got %q", out.Message)
	}
}

func TestApply_RedactsEmail(t *testing.T) {
	pattern := `[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}`
	r, _ := New([]string{pattern}, "[EMAIL]")
	e := makeEntry("user logged in as alice@example.com today")
	out := r.Apply(e)
	expected := "user logged in as [EMAIL] today"
	if out.Message != expected {
		t.Errorf("expected %q, got %q", expected, out.Message)
	}
}

func TestApply_MultiplePatterns(t *testing.T) {
	r, _ := New([]string{`\d{4}-\d{4}-\d{4}-\d{4}`, `tok_[A-Za-z0-9]+`}, "[REDACTED]")
	e := makeEntry("card 1234-5678-9012-3456 token tok_abcXYZ123")
	out := r.Apply(e)
	if out.Message != "card [REDACTED] token [REDACTED]" {
		t.Errorf("unexpected message: %q", out.Message)
	}
}

func TestApply_OriginalUnmutated(t *testing.T) {
	r, _ := New([]string{`secret`}, "[REDACTED]")
	e := makeEntry("my secret value")
	r.Apply(e)
	if e.Message != "my secret value" {
		t.Error("original entry was mutated")
	}
}

func TestApplyAll_ReturnsCorrectCount(t *testing.T) {
	r, _ := New([]string{`\d+`}, "NUM")
	entries := []parser.Entry{makeEntry("line 1"), makeEntry("line 2"), makeEntry("no digits here")}
	result := r.ApplyAll(entries)
	if len(result) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(result))
	}
	if result[0].Message != "line NUM" {
		t.Errorf("unexpected: %q", result[0].Message)
	}
}
