package approxity

import "io"

// Filter describes AMQ filter interface.
type Filter[T Hashable] interface {
	io.ReaderFrom
	io.WriterTo
	// Set add new key to the filter.
	Set(key T) error
	// Unset remove key from the filter.
	Unset(key T) error
	// Contains check if key is in the filter.
	Contains(key T) bool
	// HSet add new precalculated hash key to the filter.
	HSet(hkey uint64) error
	// HUnset remove precalculated hash key from the filter.
	HUnset(hkey uint64) error
	// HContains check if precalculated hash key is in the filter.
	HContains(hkey uint64) bool
	// Capacity returns filter capacity.
	Capacity() uint64
	// Size returns number of items added to the filter.
	Size() uint64
	// Reset flushes the filter.
	Reset()
}
