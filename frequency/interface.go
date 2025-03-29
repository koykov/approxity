package frequency

import (
	"io"

	"github.com/koykov/approxity"
)

type base[T approxity.Hashable] interface {
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
	// Reset flushes the counter.
	Reset()
}

type Estimator[T approxity.Hashable] interface {
	base[T]
	// Estimate returns frequency estimation of key.
	Estimate(key T) uint64
	// HEstimate returns frequency estimation of precalculated hash key.
	HEstimate(hkey uint64) uint64
}

type SignedEstimator[T approxity.Hashable] interface {
	base[T]
	// Estimate returns frequency estimation of key.
	Estimate(key T) int64
	// HEstimate returns frequency estimation of precalculated hash key.
	HEstimate(hkey uint64) int64
}
