package approxity

import "errors"

var (
	ErrInvalidConfig = errors.New("invalid or empty config")
	ErrNoItemsNumber = errors.New("desired number of items must be greater than 0")
	ErrInvalidFPP    = errors.New("false positive probability must be in range (0..1]")
	ErrNoHasher      = errors.New("no hasher provided")
	ErrUnsupportedOp = errors.New("unsupported operation")
	ErrEncoding      = errors.New("item encoding error")
)
