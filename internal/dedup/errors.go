package dedup

import "errors"

// ErrInvalidWindow is returned when windowSize is less than 1.
var ErrInvalidWindow = errors.New("dedup: window size must be >= 1")
