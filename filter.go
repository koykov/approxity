package bloom

import (
	"sync"

	"github.com/koykov/bitvector"
)

type Filter struct {
	once   sync.Once
	conf   *Config
	vec    bitvector.Interface
	access Access

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
	f.once.Do(f.init)
	if f.err != nil {
		return f.err
	}
	return nil
}

func (f *Filter) Check(key any) bool {
	f.once.Do(f.init)
	if f.err != nil {
		return false
	}
	return false
}

func (f *Filter) SetAccess(access Access) error {
	f.once.Do(f.init)
	if f.err != nil {
		return f.err
	}
	if f.conf.Policy == PolicySimultaneousReadWrite {
		return ErrSetAccess
	}
	f.access = access
	return nil
}

func (f *Filter) init() {
	switch f.conf.Policy {
	case PolicySimultaneousReadWrite:
		f.vec, f.err = bitvector.NewConcurrentVector(f.conf.Size, f.conf.WriteLimit)
		f.access = AccessReadWrite
	case PolicyExclusiveReadOrWrite:
		f.vec, f.err = bitvector.NewVector(f.conf.Size)
		f.access = AccessWrite
	default:
		f.err = ErrBadPolicy
	}
	if f.err != nil {
		return
	}

	if f.conf.Hasher != nil {
		f.err = ErrNoHasher
		return
	}
}
