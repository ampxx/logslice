// Package redact provides a Redactor type that masks sensitive substrings
// inside log entry messages using regular-expression patterns.
//
// # Overview
//
// Log files often contain personally identifiable information (PII) such as
// email addresses, IP addresses, or authentication tokens. The Redactor
// replaces any match of a caller-supplied set of patterns with a fixed
// placeholder string (e.g. "[REDACTED]"), leaving the rest of the entry —
// including the timestamp — intact.
//
// # Usage
//
//	r, err := redact.New(
//		[]string{`\b[\w.+-]+@[\w-]+\.[\w.]+\b`, `\b(?:\d{1,3}\.){3}\d{1,3}\b`},
//		"[REDACTED]",
//	)
//	if err != nil {
//		log.Fatal(err)
//	}
//	cleaned := r.Apply(entry)
package redact
