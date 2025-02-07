package bloom

import (
	"sync"

	"github.com/koykov/bitvector"
)

type Filter struct {
	once sync.Once
	conf *Config
	vec  bitvector.Interface

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
	return nil
}

func (f *Filter) Check(key any) bool {
	if f.once.Do(f.init); f.err != nil {
		return false
	}
	return false
}

func (f *Filter) init() {
	if f.conf.Hasher != nil {
		f.err = ErrNoHasher
		return
	}
	if f.vec, f.err = bitvector.NewVector(f.conf.Size); f.err != nil {
		return
	}
}
