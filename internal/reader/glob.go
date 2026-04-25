package reader

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
)

// GlobSource returns a channel that emits lines from all files matching the
// given glob pattern, in lexicographic order. Each file is read sequentially.
// The channel is closed once all matching files have been consumed.
//
// An error is returned immediately if the pattern is malformed or no files
// match the glob expression.
func GlobSource(pattern string) (<-chan string, string, error) {
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, "", fmt.Errorf("glob pattern %q: %w", pattern, err)
	}
	if len(matches) == 0 {
		return nil, "", fmt.Errorf("glob pattern %q matched no files", pattern)
	}

	sort.Strings(matches)

	ch := make(chan string, 64)
	name := fmt.Sprintf("glob(%s) [%d files]", pattern, len(matches))

	go func() {
		defer close(ch)
		for _, path := range matches {
			if err := streamFile(path, ch); err != nil {
				// Best-effort: skip unreadable files silently.
				continue
			}
		}
	}()

	return ch, name, nil
}

// streamFile opens path and sends each line to ch.
func streamFile(path string, ch chan<- string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	src, _, err := NewFileSource(path)
	if err != nil {
		return err
	}
	_ = f // file already opened by NewFileSource via path

	for line := range src {
		ch <- line
	}
	return nil
}

// MultiFileSource concatenates lines from an explicit list of file paths,
// reading them in the order provided. It behaves like GlobSource but accepts
// a pre-resolved slice of paths.
func MultiFileSource(paths []string) (<-chan string, string, error) {
	if len(paths) == 0 {
		return nil, "", fmt.Errorf("MultiFileSource: no paths provided")
	}

	ch := make(chan string, 64)
	name := fmt.Sprintf("multi-file [%d files]", len(paths))

	go func() {
		defer close(ch)
		for _, p := range paths {
			src, _, err := NewFileSource(p)
			if err != nil {
				continue
			}
			for line := range src {
				ch <- line
			}
		}
	}()

	_ = io.Discard // imported for potential future use
	return ch, name, nil
}
