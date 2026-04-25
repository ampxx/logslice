// Package highlight provides ANSI terminal colour highlighting for
// substrings within log output.
//
// # Overview
//
// A [Highlighter] is constructed with a regular-expression pattern and
// an ANSI colour code.  Its [Highlighter.Apply] method returns a copy
// of the input line with every match wrapped in the chosen colour
// sequence, leaving the rest of the line unchanged.
//
// # Usage
//
//	h, err := highlight.New("ERROR|WARN", highlight.Red)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(h.Apply(line))
//
// Predefined colour constants (Red, Yellow, Cyan, Bold) are provided
// for convenience.  Pass an empty pattern to obtain a no-op
// Highlighter that returns lines unmodified.
package highlight
