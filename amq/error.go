package amq

import "errors"

var (
	ErrNoItemsNumber = errors.New("desired number of items must be greater than 0")
	ErrInvalidFPP    = errors.New("false positive probability must be in range (0..1]")
)
