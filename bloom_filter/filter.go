package bloom

import (
	"fmt"
	"io"
	"sync"

	"github.com/koykov/approxity"
	"github.com/koykov/bitvector"
)

// Bloom filter implementation.
// By default, filter doesn't support concurrent read/write operations - you must set up the filter before reading.
// Concurrent reading allowed afterward.
// If you want to use concurrent read/write operations, fill up Concurrent section in Config object.
type filter[T approxity.Hashable] struct {
	approxity.Base[T]
	once sync.Once
	conf *Config
	m, k uint64
	vec  bitvector.Interface

	err error
}

// NewFilter creates new filter.
func NewFilter[T approxity.Hashable](config *Config) (approxity.Filter[T], error) {
	if config == nil {
		return nil, approxity.ErrInvalidConfig
	}
	f := &filter[T]{
		conf: config.copy(),
	}
	if f.once.Do(f.init); f.err != nil {
		return nil, f.err
	}
	return f, nil
}

// Set adds new key to the filter.
func (f *filter[T]) Set(key T) error {
	if f.once.Do(f.init); f.err != nil {
		return f.err
	}
	for i := uint64(0); i < f.k; i++ {
		h, err := f.h(key, i)
		if err != nil {
			return f.mw().Set(err)
		}
		f.vec.Set(h % f.m)
	}
	return f.mw().Set(nil)
}

// HSet sets new predefined hash key to the filter.
func (f *filter[T]) HSet(hkey uint64) error {
	f.vec.Set(hkey % f.m)
	return f.mw().Set(nil)
}

// Unset removes key from the filter.
// Caution! Bloom filter doesn't support this operation!
func (f *filter[T]) Unset(_ T) error {
	return f.mw().Unset(approxity.ErrUnsupportedOp)
}

// HUnset removes predefined hash key from the filter.
// Caution! Bloom filter doesn't support this operation!
func (f *filter[T]) HUnset(_ uint64) error {
	return f.mw().Unset(approxity.ErrUnsupportedOp)
}

// Contains checks if key is in the filter.
func (f *filter[T]) Contains(key T) bool {
	if f.once.Do(f.init); f.err != nil {
		return false
	}
	for i := uint64(0); i < f.k; i++ {
		h, err := f.h(key, i)
		if err != nil {
			return f.mw().Contains(false)
		}
		if f.vec.Get(h%f.m) == 0 {
			return f.mw().Contains(false)
		}
	}
	return f.mw().Contains(true)
}

// HContains checks if predefined hash key is in the filter.
func (f *filter[T]) HContains(hkey uint64) bool {
	return f.mw().Contains(f.vec.Get(hkey%f.m) == 1)
}

// Capacity returns filter capacity.
func (f *filter[T]) Capacity() uint64 {
	return f.vec.Capacity()
}

// Size returns number of items added to the filter.
func (f *filter[T]) Size() uint64 {
	return f.vec.Size()
}

func (f *filter[T]) ReadFrom(r io.Reader) (int64, error) {
	if f.once.Do(f.init); f.err != nil {
		return 0, f.err
	}
	expect := f.vec.Capacity() / 8 // bitvector returns capacity in bits, so recalculate to bytes
	n, err := f.vec.ReadFrom(r)
	if err != nil {
		return n, err
	}
	hsz := uint64(32) // header size of vector in bytes
	if f.conf.Concurrent != nil {
		hsz = 40 // header size of concurrent vector
	}
	if actual := uint64(n) - hsz; actual != expect {
		return n, fmt.Errorf("expected %d bytes, but got %d", expect, actual)
	}
	return n, nil
}

func (f *filter[T]) WriteTo(w io.Writer) (int64, error) {
	if f.once.Do(f.init); f.err != nil {
		return 0, f.err
	}
	return f.vec.WriteTo(w)
}

// Reset flushes filter data.
func (f *filter[T]) Reset() {
	if f.once.Do(f.init); f.err != nil {
		return
	}
	f.vec.Reset()
	f.mw().Reset()
}

func (f *filter[T]) init() {
	c := f.conf
	if c.ItemsNumber == 0 {
		f.err = approxity.ErrNoItemsNumber
		return
	}
	if c.Hasher == nil {
		f.err = approxity.ErrNoHasher
		return
	}
	if c.MetricsWriter == nil {
		c.MetricsWriter = approxity.DummyMetricsWriter{}
	}
	if c.FPP == 0 {
		c.FPP = defaultFPP
	}
	if c.FPP < 0 || c.FPP > 1 {
		f.err = approxity.ErrInvalidFPP
		return
	}

	f.m = optimalM(c.ItemsNumber, c.FPP)
	f.k = optimalK(c.ItemsNumber, f.m)
	if c.Concurrent != nil {
		f.vec, f.err = bitvector.NewConcurrentVector(f.m, c.Concurrent.WriteAttemptsLimit)
	} else {
		f.vec, f.err = bitvector.NewVector(f.m)
	}
	f.mw().Capacity(f.m)
}

func (f *filter[T]) h(key T, salt uint64) (uint64, error) {
	return f.HashSalt(f.conf.Hasher, key, salt)
}

func (f *filter[T]) mw() approxity.MetricsWriter {
	return f.conf.MetricsWriter
}
