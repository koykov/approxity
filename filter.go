package bloom

import (
	"sync"

	"github.com/koykov/bitvector"
)

type Filter struct {
	once sync.Once
	conf *Config
	vec  bitvector.Interface
	cncr bool

	err error
}

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
	if f.cncr {
		f.vec, f.err = bitvector.NewConcurrentVector(f.conf.Size, f.conf.Concurrent.WriteAttemptsLimit)
	} else {
		f.vec, f.err = bitvector.NewVector(f.conf.Size)
	}
}
