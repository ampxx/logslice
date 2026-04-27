// Package tail provides live log-following functionality, similar to
// `tail -f`, for use with the logslice pipeline.
package tail

import (
	"bufio"
	"context"
	"io"
	"os"
	"time"
)

// DefaultPollInterval is how often the tailer checks for new content.
const DefaultPollInterval = 250 * time.Millisecond

// Tailer streams new lines appended to a file, blocking until the context
// is cancelled or an unrecoverable error occurs.
type Tailer struct {
	path         string
	pollInterval time.Duration
}

// New returns a Tailer for the given file path.
func New(path string) *Tailer {
	return &Tailer{path: path, pollInterval: DefaultPollInterval}
}

// WithPollInterval overrides the default polling interval.
func (t *Tailer) WithPollInterval(d time.Duration) *Tailer {
	t.pollInterval = d
	return t
}

// Lines sends each newly appended line to the returned channel.
// The channel is closed when ctx is cancelled or a read error occurs.
func (t *Tailer) Lines(ctx context.Context) (<-chan string, error) {
	f, err := os.Open(t.path)
	if err != nil {
		return nil, err
	}

	// Seek to end so we only tail new content.
	if _, err := f.Seek(0, io.SeekEnd); err != nil {
		f.Close()
		return nil, err
	}

	ch := make(chan string, 64)
	go func() {
		defer close(ch)
		defer f.Close()

		reader := bufio.NewReader(f)
		ticker := time.NewTicker(t.pollInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				for {
					line, err := reader.ReadString('\n')
					if len(line) > 0 {
						// Strip trailing newline before sending.
						if len(line) > 0 && line[len(line)-1] == '\n' {
							line = line[:len(line)-1]
						}
						select {
						case ch <- line:
						case <-ctx.Done():
							return
						}
					}
					if err != nil {
						// io.EOF is normal — no new data yet.
						break
					}
				}
			}
		}
	}()

	return ch, nil
}
