package amq

import "errors"

var (
	ErrBadConfig     = errors.New("invalid or empty config")
	ErrNoItemsNumber = errors.New("desired number of items must be greater than 0")
	ErrNoHasher      = errors.New("no hasher provided")
	ErrUnsupportedOp = errors.New("unsupported operation")
)
