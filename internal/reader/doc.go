// Package reader abstracts log input sources for logslice.
//
// It supports reading from regular files and from standard input,
// exposing a uniform channel-based Lines() API that integrates
// cleanly with the parser and filter pipeline.
//
// Typical usage:
//
//	src, err := reader.NewFileSource("app.log")
//	if err != nil { ... }
//	defer src.Close()
//
//	for line := range src.Lines() {
//		// pass line to parser.LineParser
//	}
//
// For stdin:
//
//	src := reader.NewStdinSource()
//	for line := range src.Lines() { ... }
package reader
