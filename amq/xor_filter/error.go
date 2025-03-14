package xor

import "errors"

var (
	ErrUnsupportedSet   = errors.New("filter doesn't support setting new items, create new filter with new keys list instead")
	ErrUnsupportedUnset = errors.New("filter doesn't support items deletion, create new filter with new keys list instead")
	ErrEmptyKeyset      = errors.New("keys list is empty")
)
