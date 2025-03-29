package pbtk

import "errors"

var (
	ErrInvalidConfig    = errors.New("invalid or empty config")
	ErrNoHasher         = errors.New("no hasher provided")
	ErrUnsupportedOp    = errors.New("unsupported operation")
	ErrEncoding         = errors.New("item encoding error")
	ErrWriteLimitExceed = errors.New("write limit exceeded")
	ErrInvalidSignature = errors.New("invalid signature")
	ErrVersionMismatch  = errors.New("version mismatch")
)
