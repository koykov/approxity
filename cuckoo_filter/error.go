package cuckoo

import "errors"

var (
	ErrFullBucket      = errors.New("bucket is full")
	ErrFullFilter      = errors.New("filter is full")
	ErrWriteLimitReach = errors.New("write limit reached")
)
