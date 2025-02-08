package bloom

import (
	"sync"

	"github.com/koykov/bitvector"
)

// Filter represents Bloom filter.
// By default, filter doesn't support concurrent read/write operations - you must set up the filter before reading.
// Concurrent reading allowed afterward.
// If you want to use concurrent read/write operations, fill up Concurrent section in Config object.
type Filter struct {
	once sync.Once
	conf *Config
	vec  bitvector.Interface

	err error
}

// NewFilter creates new Bloom filter.
func NewFilter(config *Config) (*Filter, error) {
	if config == nil {
		return nil, ErrBadConfig
	}
	f := &Filter{
		conf: config.copy(),
	}
	f.once.Do(f.init)
	return f, f.err
}

// Set adds new item to the filter.
func (f *Filter) Set(key any) error {
	if f.once.Do(f.init); f.err != nil {
		return f.err
	}
	for i := uint64(0); i < f.conf.HashChecksLimit+1; i++ {
		h := f.conf.Hasher.Hash(key) + i
		f.vec.Set(h % f.conf.Size)
	}
	return nil
}

// Clear removes item from the filter.
func (f *Filter) Clear(key any) error {
	if f.once.Do(f.init); f.err != nil {
		return f.err
	}
	for i := uint64(0); i < f.conf.HashChecksLimit+1; i++ {
		h := f.conf.Hasher.Hash(key) + i
		f.vec.Clear(h % f.conf.Size)
	}
	return nil
}

// Check checks if item is in the filter.
func (f *Filter) Check(key any) bool {
	if f.once.Do(f.init); f.err != nil {
		return false
	}
	for i := uint64(0); i < f.conf.HashChecksLimit+1; i++ {
		h := f.conf.Hasher.Hash(key) + i
		if f.vec.Get(h%f.conf.Size) == 1 {
			return true
		}
	}
	return false
}

func (f *Filter) init() {
	if f.conf.Hasher == nil {
		f.err = ErrNoHasher
		return
	}
	if f.conf.Concurrent != nil {
		f.vec, f.err = bitvector.NewConcurrentVector(f.conf.Size, f.conf.Concurrent.WriteAttemptsLimit)
	} else {
		f.vec, f.err = bitvector.NewVector(f.conf.Size)
	}
}
