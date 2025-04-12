package frequency

import (
	"context"
	"io"

	"github.com/koykov/pbtk"
)

type base[T pbtk.Hashable] interface {
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

type Estimator[T pbtk.Hashable] interface {
	base[T]
	// Estimate returns frequency estimation of key.
	Estimate(key T) uint64
	// HEstimate returns frequency estimation of precalculated hash key.
	HEstimate(hkey uint64) uint64
}

type SignedEstimator[T pbtk.Hashable] interface {
	base[T]
	// Estimate returns signed frequency estimation of key.
	Estimate(key T) int64
	// HEstimate returns signed frequency estimation of precalculated hash key.
	HEstimate(hkey uint64) int64
}

type PreciseEstimator[T pbtk.Hashable] interface {
	base[T]
	// Estimate returns float frequency estimation of key.
	Estimate(key T) float64
	// HEstimate returns float frequency estimation of precalculated hash key.
	HEstimate(hkey uint64) float64
}

type Decayer interface {
	// Decay applies factor to all counters inside.
	Decay(ctx context.Context, factor float64) error
}
