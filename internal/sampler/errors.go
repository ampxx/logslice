package sampler

import "errors"

// ErrInvalidN is returned when n < 1 is supplied for ModeNth.
var ErrInvalidN = errors.New("sampler: n must be >= 1 for nth mode")

// ErrInvalidProb is returned when probability is outside (0, 1].
var ErrInvalidProb = errors.New("sampler: probability must be in range (0, 1]")
