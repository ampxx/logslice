package redact

import "errors"

// ErrEmptyPlaceholder is returned when an empty placeholder string is provided
// to New.
var ErrEmptyPlaceholder = errors.New("redact: placeholder must not be empty")
