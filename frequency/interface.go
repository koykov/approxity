package frequency

import (
	"io"

	"github.com/koykov/approxity"
)

type Estimator[T approxity.Hashable] interface {
	io.ReaderFrom
	io.WriterTo
	// Add adds new key to the counter.
	Add(key T) error
	// AddN adds new key to the counter with given count.
	AddN(key T, n uint64) error
	// HAdd adds new precalculated hash key to the counter.
	HAdd(hkey uint64) error
	// HAddN adds new precalculated hash key to the counter with given count.
	HAddN(hkey uint64, n uint64) error
	// Estimate returns frequency estimation of key.
	Estimate(key T) uint64
	// HEstimate returns frequency estimation of precalculated hash key.
	HEstimate(hkey uint64) uint64
	// Reset flushes the counter.
	Reset()
}
