package xor

import (
	"sync"

	"github.com/koykov/approxity"
	"github.com/koykov/approxity/amq"
)

var p sync.Pool

func AcquireWithKeys[T approxity.Hashable](config *Config, keys []T) (_ amq.Filter[T], err error) {
	var f amq.Filter[T]
	if v := p.Get(); v != nil {
		ff := v.(*filter[T])
		ff.conf = config.copy()
		f = ff
	} else {
		f, err = NewFilterWithKeys(config, keys)
	}
	return f, err
}

func AcquireWithHKeys(config *Config, hkeys []uint64) (_ amq.Filter[uint64], err error) {
	var f amq.Filter[uint64]
	if v := p.Get(); v != nil {
		ff := v.(*filter[uint64])
		ff.conf = config.copy()
		f = ff
	} else {
		f, err = NewFilterWithHKeys(config, hkeys)
	}
	return f, err
}

func Release[T approxity.Hashable](f amq.Filter[T]) {
	f.Reset()
	p.Put(f)
}
