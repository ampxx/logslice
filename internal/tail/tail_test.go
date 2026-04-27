package tail_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/dleviminzi/logslice/internal/tail"
)

func TestTailer_ReceivesAppendedLines(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "tail-*.log")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	tr := tail.New(f.Name()).WithPollInterval(20 * time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	ch, err := tr.Lines(ctx)
	if err != nil {
		t.Fatal(err)
	}

	// Write lines after the tailer has started.
	lines := []string{"first line", "second line", "third line"}
	go func() {
		time.Sleep(40 * time.Millisecond)
		for _, l := range lines {
			f.WriteString(l + "\n")
		}
	}()

	received := make([]string, 0, len(lines))
	for len(received) < len(lines) {
		select {
		case line, ok := <-ch:
			if !ok {
				t.Fatalf("channel closed early; got %d/%d lines", len(received), len(lines))
			}
			received = append(received, line)
		case <-ctx.Done():
			t.Fatalf("timeout; got %d/%d lines", len(received), len(lines))
		}
	}

	for i, want := range lines {
		if received[i] != want {
			t.Errorf("line %d: want %q, got %q", i, want, received[i])
		}
	}
}

func TestTailer_ClosesChannelOnCancel(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "tail-cancel-*.log")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()

	ctx, cancel := context.WithCancel(context.Background())
	tr := tail.New(f.Name()).WithPollInterval(20 * time.Millisecond)

	ch, err := tr.Lines(ctx)
	if err != nil {
		t.Fatal(err)
	}

	cancel()

	select {
	case _, ok := <-ch:
		if ok {
			t.Error("expected channel to be closed after cancel")
		}
	case <-time.After(time.Second):
		t.Error("channel was not closed within 1 s after cancel")
	}
}

func TestTailer_MissingFile(t *testing.T) {
	tr := tail.New("/nonexistent/path/logfile.log")
	ctx := context.Background()
	_, err := tr.Lines(ctx)
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}
