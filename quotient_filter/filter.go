package quotient

import (
	"io"
	"sync"

	"github.com/koykov/amq"
)

type Filter struct {
	conf       *Config
	once       sync.Once
	qb, rb     uint64 // quotient and remainder bits
	bs         uint64 // bucket size (rb+3)
	m          uint64 // total filter size
	bm, qm, rm uint64 // bucket mask, quotient mask, remainder mask
	vec        []uint64

	err error
}

func NewFilter(config *Config) (*Filter, error) {
	if config == nil {
		return nil, amq.ErrInvalidConfig
	}
	f := &Filter{
		conf: config.copy(),
	}
	if f.once.Do(f.init); f.err != nil {
		return nil, f.err
	}
	return f, nil
}

func (f *Filter) Set(key any) error {
	// todo implement me
	return nil
}

func (f *Filter) HSet(hkey uint64) error {
	// todo implement me
	return nil
}

func (f *Filter) Unset(key any) error {
	// todo implement me
	return nil
}

func (f *Filter) HUnset(hkey uint64) error {
	// todo implement me
	return nil
}

func (f *Filter) Contains(key any) bool {
	// todo implement me
	return false
}

func (f *Filter) HContains(hkey uint64) bool {
	// todo implement me
	return false
}

func (f *Filter) Reset() {
	// todo implement me
}

func (f *Filter) ReadFrom(r io.Reader) (int64, error) {
	// todo implement me
	return 0, nil
}

func (f *Filter) WriteTo(w io.Writer) (int64, error) {
	// todo implement me
	return 0, nil
}

func (f *Filter) init() {
	c := f.conf
	if c.ItemsNumber == 0 {
		f.err = amq.ErrNoItemsNumber
		return
	}
	if c.Hasher == nil {
		f.err = amq.ErrNoHasher
		return
	}
	if c.MetricsWriter == nil {
		c.MetricsWriter = amq.DummyMetricsWriter{}
	}
	if c.FPP == 0 {
		c.FPP = defaultFPP
	}
	if c.FPP < 0 || c.FPP > 1 {
		f.err = amq.ErrInvalidFPP
		return
	}
	if c.LoadFactor == 0 {
		c.LoadFactor = defaultLoadFactor
	}
	if c.LoadFactor < 0 || c.LoadFactor > 1 {
		f.err = ErrInvalidLoadFactor
		return
	}

	if f.m, f.qb, f.rb = optimalMQR(c.ItemsNumber, c.FPP, c.LoadFactor); f.qb+f.qb > 64 {
		f.err = ErrBucketOverflow
		return
	}
	f.bs = f.rb + 3
	f.vec = make([]uint64, f.m)
	f.mw().Capacity(f.m)

	f.qm, f.rm, f.bm = lowMask(f.qb), lowMask(f.rb), lowMask(f.bs)
}

func (f *Filter) mw() amq.MetricsWriter {
	return f.conf.MetricsWriter
}

func lowMask(v uint64) uint64 {
	return (1 << v) - 1
}
