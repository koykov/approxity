package countminsketch

import "errors"

var (
	ErrInvalidConfidence = errors.New("confidence must be in range (0..1)")
	ErrInvalidEpsilon    = errors.New("epsilon must be in range (0..1)")
)
