package bloom

import (
	"sync"

	"github.com/koykov/bitvector"
)

type ConcurrentFilter struct {
	once sync.Once
	conf *Config
	vec  *bitvector.ConcurrentVector

	err error
}

func NewConcurrentFilter(config *Config) (*ConcurrentFilter, error) {
	if config == nil {
		return nil, ErrBadConfig
	}
	f := &ConcurrentFilter{
		conf: config.copy(),
	}
	f.once.Do(f.init)
	return f, nil
}

func (f *ConcurrentFilter) Set(key any) error {
	if f.once.Do(f.init); f.err != nil {
		return f.err
	}
	for i := uint64(0); i < f.conf.HashChecksLimit+1; i++ {
		h := f.conf.Hasher.Hash(key) + i
		f.vec.Set(h % f.conf.Size)
	}
	return nil
}

func (f *ConcurrentFilter) Check(key any) bool {
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

func (f *ConcurrentFilter) init() {
	if f.conf.Hasher == nil {
		f.err = ErrNoHasher
		return
	}
	if f.vec, f.err = bitvector.NewConcurrentVector(f.conf.Size, f.conf.Concurrent.WriteAttemptsLimit); f.err != nil {
		return
	}
}
