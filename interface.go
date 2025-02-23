package amq

// Interface describes AMQ struct interface.
type Interface interface {
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
	// Reset flushes the filter.
	Reset()
}
