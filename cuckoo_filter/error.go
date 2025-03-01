package cuckoo

import "errors"

var (
	ErrFullBucket       = errors.New("bucket is full")
	ErrFullFilter       = errors.New("filter is full")
	ErrWriteLimitReach  = errors.New("write limit reached")
	ErrInvalidSignature = errors.New("invalid vector signature")
	ErrVersionMismatch  = errors.New("vector version mismatch")
)
