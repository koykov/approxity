package bloom

import (
	"sync"

	"github.com/koykov/amq"
	"github.com/koykov/bitvector"
)

// Filter represents Bloom filter.
// By default, filter doesn't support concurrent read/write operations - you must set up the filter before reading.
// Concurrent reading allowed afterward.
// If you want to use concurrent read/write operations, fill up Concurrent section in Config object.
type Filter struct {
	amq.Base
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

// Set adds new key to the filter.
func (f *Filter) Set(key any) error {
	if f.once.Do(f.init); f.err != nil {
		return f.err
	}
	for i := uint64(0); i < f.c().HashChecksLimit+1; i++ {
		h, err := f.h(key, i)
		if err != nil {
			return f.mw().Set(err)
		}
		f.vec.Set(h % f.c().Size)
	}
	return f.mw().Set(nil)
}

// Unset removes key from the filter.
func (f *Filter) Unset(key any) error {
	if f.once.Do(f.init); f.err != nil {
		return f.err
	}
	for i := uint64(0); i < f.c().HashChecksLimit+1; i++ {
		h, err := f.h(key, i)
		if err != nil {
			return f.mw().Unset(err)
		}
		f.vec.Unset(h % f.c().Size)
	}
	return f.mw().Unset(nil)
}

// Contains checks if key is in the filter.
func (f *Filter) Contains(key any) bool {
	if f.once.Do(f.init); f.err != nil {
		return false
	}
	for i := uint64(0); i < f.c().HashChecksLimit+1; i++ {
		h, err := f.h(key, i)
		if err != nil {
			return f.mw().Contains(false)
		}
		if f.vec.Get(h%f.c().Size) == 0 {
			return f.mw().Contains(false)
		}
	}
	return f.mw().Contains(true)
}

func (f *Filter) Reset() {
	if f.once.Do(f.init); f.err != nil {
		return
	}
	f.vec.Reset()
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
	if f.conf.MetricsWriter == nil {
		f.conf.MetricsWriter = amq.DummyMetricsWriter{}
	}
}

func (f *Filter) c() *Config {
	return f.conf
}

func (f *Filter) h(key any, seed uint64) (uint64, error) {
	return f.Hash(f.c().Hasher, key, seed)
}

func (f *Filter) mw() amq.MetricsWriter {
	return f.conf.MetricsWriter
}
