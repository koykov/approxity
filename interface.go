package amq

import "io"

// Filter describes AMQ filter interface.
type Filter interface {
	io.ReaderFrom
	io.WriterTo
	// Set add new key to the filter.
	Set(key any) error
	// Unset remove key from the filter.
	Unset(key any) error
	// Contains check if key is in the filter.
	Contains(key any) bool
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

// Counter describes AMQ unique counter interface.
type Counter interface {
	io.ReaderFrom
	io.WriterTo
	// Add adds new key to the counter.
	Add(key any) error
	// HAdd adds new precalculated hash key to the counter.
	HAdd(hkey uint64) error
	// Count returns number of unique keys added to the counter.
	Count() uint64
}
