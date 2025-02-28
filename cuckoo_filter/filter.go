package cuckoo

import (
	"math/bits"
	"math/rand"
	"sync"

	"github.com/koykov/amq"
)

const bucketsz = 4

type Filter struct {
	amq.Base
	once sync.Once
	conf *Config

	vec ivector
	m   uint64
	bp  uint64
	hsh [256]uint64

	err error
}

func NewFilter(conf *Config) (*Filter, error) {
	f := &Filter{
		conf: conf.copy(),
	}
	if f.once.Do(f.init); f.err != nil {
		return nil, f.err
	}
	return f, nil
}

func (f *Filter) Set(key any) error {
	if f.once.Do(f.init); f.err != nil {
		return f.mw().Set(f.err)
	}
	i0, i1, fp, err := f.calcI2FP(key, f.bp, 0)
	if err != nil {
		return f.mw().Set(err)
	}
	return f.hset(i0, i1, fp)
}

func (f *Filter) HSet(hkey uint64) error {
	if f.once.Do(f.init); f.err != nil {
		return f.mw().Set(f.err)
	}
	i0, i1, fp, err := f.hcalcI2FP(hkey, f.bp)
	if err != nil {
		return f.mw().Set(err)
	}
	return f.hset(i0, i1, fp)
}

func (f *Filter) hset(i0, i1 uint64, fp byte) (err error) {
	if err = f.vec.add(i0, fp); err == nil {
		return f.mw().Set(nil)
	}
	if err = f.vec.add(i1, fp); err == nil {
		return f.mw().Set(nil)
	}
	i := i0
	if rand.Intn(2) == 1 {
		i = i1
	}
	for k := uint64(0); k < f.c().KicksLimit; k++ {
		j := uint64(rand.Intn(bucketsz))
		pfp := fp
		fp = f.vec.fpv(i, j)
		_ = f.vec.set(i, j, pfp)

		m := mask64[f.bp]
		i = (i & m) ^ (f.hsh[fp] & m)
		if err = f.vec.add(i, fp); err == nil {
			return f.mw().Set(nil)
		}
	}
	return f.mw().Set(ErrFullFilter)
}

func (f *Filter) Unset(key any) error {
	if f.once.Do(f.init); f.err != nil {
		return f.mw().Unset(f.err)
	}
	i0, i1, fp, err := f.calcI2FP(key, f.bp, 0)
	if err != nil {
		return f.mw().Unset(err)
	}
	return f.hunset(i0, i1, fp)
}

func (f *Filter) HUnset(hkey uint64) error {
	if f.once.Do(f.init); f.err != nil {
		return f.mw().Unset(f.err)
	}
	i0, i1, fp, err := f.hcalcI2FP(hkey, f.bp)
	if err != nil {
		return f.mw().Unset(err)
	}
	return f.hunset(i0, i1, fp)
}

func (f *Filter) hunset(i0, i1 uint64, fp byte) (err error) {
	if f.vec.unset(i0, fp) {
		return f.mw().Unset(nil)
	}
	f.vec.unset(i1, fp)
	return f.mw().Unset(nil)
}

func (f *Filter) Contains(key any) bool {
	if f.once.Do(f.init); f.err != nil {
		return f.mw().Contains(false)
	}
	i0, i1, fp, err := f.calcI2FP(key, f.bp, 0)
	if err != nil {
		return f.mw().Contains(false)
	}
	return f.hcontains(i0, i1, fp)
}

func (f *Filter) HContains(hkey uint64) bool {
	if f.once.Do(f.init); f.err != nil {
		return f.mw().Contains(false)
	}
	i0, i1, fp, err := f.hcalcI2FP(hkey, f.bp)
	if err != nil {
		return f.mw().Contains(false)
	}
	return f.hcontains(i0, i1, fp)
}

func (f *Filter) hcontains(i0, i1 uint64, fp byte) bool {
	if f.vec.fpi(i0, fp) != -1 {
		return f.mw().Contains(true)
	}
	return f.mw().Contains(f.vec.fpi(i1, fp) != -1)
}

func (f *Filter) Capacity() uint64 {
	return f.vec.capacity()
}

func (f *Filter) Size() uint64 {
	return f.vec.size()
}

func (f *Filter) Reset() {
	if f.once.Do(f.init); f.err != nil {
		return
	}
	f.vec.reset()
	f.mw().Reset()
}

func (f *Filter) calcI2FP(key any, bp, i uint64) (i0 uint64, i1 uint64, fp byte, err error) {
	var hkey uint64
	if hkey, err = f.Hash(f.c().Hasher, key); err != nil {
		return
	}
	return f.hcalcI2FP(hkey, bp)
}

func (f *Filter) hcalcI2FP(hkey, bp uint64) (i0, i1 uint64, fp byte, err error) {
	fp = byte(hkey%255 + 1)
	i0 = (hkey >> 32) & mask64[bp]
	m := mask64[bp]
	i1 = (i0 & m) ^ (f.hsh[fp] & m)
	return
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
	if c.KicksLimit == 0 {
		c.KicksLimit = defaultKicksLimit
	}

	f.m = optimalM(c.ItemsNumber)
	f.bp = uint64(bits.TrailingZeros64(f.m))
	if f.m == 0 {
		f.m = 1
	}
	if c.Concurrent != nil {
		f.vec = newCnvector(f.m, c.Concurrent.WriteAttemptsLimit)
	} else {
		f.vec = newVector(f.m)
	}
	f.mw().Capacity(c.ItemsNumber)

	var buf []byte
	for i := 0; i < 256; i++ {
		buf = append(buf[:0], byte(i))
		f.hsh[i], _ = f.Hash(c.Hasher, buf)
	}
}

func (f *Filter) c() *Config {
	return f.conf
}

func (f *Filter) mw() amq.MetricsWriter {
	return f.conf.MetricsWriter
}
