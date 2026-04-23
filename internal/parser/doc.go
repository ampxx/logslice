// Package parser implements log line parsing for logslice.
//
// It provides:
//   - Entry: a structured representation of a single log line containing
//     the original raw text, an extracted timestamp, and a log level.
//   - LineParser: reads an io.Reader line-by-line and converts each
//     non-empty line into an Entry.
//   - ParseTimestamp: tries a set of CommonLayouts to convert a raw
//     timestamp string into a time.Time value.
//
// Usage:
//
//	f, _ := os.Open("app.log")
//	defer f.Close()
//
//	p := parser.NewLineParser(f)
//	entries, err := p.Parse()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	for _, e := range entries {
//	    fmt.Println(e)
//	}
package parser
