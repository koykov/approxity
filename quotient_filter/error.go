package quotient

import "errors"

var (
	ErrInvalidLoadFactor = errors.New("load factor must be in range (0..1]")
	ErrBucketOverflow    = errors.New("bucket overflow")
)
