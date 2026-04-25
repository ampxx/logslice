// Package reader provides utilities for reading log files and stdin,
// yielding raw lines for downstream parsing.
package reader

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// Source represents an input source for log lines.
type Source struct {
	Name   string
	reader io.ReadCloser
}

// NewFileSource opens a file and returns a Source for it.
func NewFileSource(path string) (*Source, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("reader: open %q: %w", path, err)
	}
	return &Source{Name: path, reader: f}, nil
}

// NewStdinSource returns a Source backed by os.Stdin.
func NewStdinSource() *Source {
	return &Source{Name: "<stdin>", reader: io.NopCloser(os.Stdin)}
}

// Lines returns a channel that emits each line from the source.
// The channel is closed when the source is exhausted or an error occurs.
// The caller should invoke Close after consuming the channel.
func (s *Source) Lines() <-chan string {
	ch := make(chan string)
	go func() {
		defer close(ch)
		scanner := bufio.NewScanner(s.reader)
		for scanner.Scan() {
			ch <- scanner.Text()
		}
	}()
	return ch
}

// Close releases the underlying resource.
func (s *Source) Close() error {
	return s.reader.Close()
}
