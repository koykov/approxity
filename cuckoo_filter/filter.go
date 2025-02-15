package cuckoo

import (
	"sync"

	"github.com/koykov/amq"
)

type Filter struct {
	amq.Base
	once sync.Once
	conf *Config

	buckets []bucket
	buf     []byte

	err error
}

func NewFilter(conf *Config) (*Filter, error) {
	f := &Filter{
		conf: conf,
	}
	f.once.Do(f.init)
	return f, f.err
}

func (f *Filter) Set(key any) error {
	// todo implement me
	return nil
}

func (f *Filter) Unset(key any) error {
	// todo implement me
	return nil
}

func (f *Filter) Contains(key any) bool {
	// todo implement me
	return false
}

func (f *Filter) Reset() {
	// todo implement me
}

func (f *Filter) init() {
	// todo implement me
}
