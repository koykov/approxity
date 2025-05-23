package xor

import (
	"sync"

	"github.com/koykov/pbtk"
	"github.com/koykov/pbtk/amq"
)

var p sync.Pool

func AcquireWithKeys[T pbtk.Hashable](config *Config, keys []T) (_ amq.Filter[T], err error) {
	if keys = pbtk.Deduplicate(keys); len(keys) == 0 {
		return nil, ErrEmptyKeyset
	}
	var f amq.Filter[T]
	if v := p.Get(); v != nil {
		ff := v.(*filter[T])
		ff.conf = config.copy()
		ff.len = uint64(len(keys))
		ff.once.Do(ff.init)
		f = ff
	} else {
		f, err = NewFilterWithKeys(config, keys)
	}
	return f, err
}

func AcquireWithHKeys(config *Config, hkeys []uint64) (_ amq.Filter[uint64], err error) {
	if hkeys = pbtk.Deduplicate(hkeys); len(hkeys) == 0 {
		return nil, ErrEmptyKeyset
	}
	var f amq.Filter[uint64]
	if v := p.Get(); v != nil {
		ff := v.(*filter[uint64])
		ff.conf = config.copy()
		ff.len = uint64(len(hkeys))
		f = ff
	} else {
		f, err = NewFilterWithHKeys(config, hkeys)
	}
	return f, err
}

func Release[T pbtk.Hashable](f amq.Filter[T]) {
	f.Reset()
	p.Put(f)
}

var _ = AcquireWithHKeys
