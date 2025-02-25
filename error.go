package amq

import "errors"

var (
	ErrBadSize       = errors.New("size must be greater than 0")
	ErrBadConfig     = errors.New("invalid or empty config")
	ErrNoHasher      = errors.New("no hasher provided")
	ErrUnsupportedOp = errors.New("unsupported operation")
)
