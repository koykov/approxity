package cuckoo

import (
	"math/bits"
	"math/rand"
	"sync"
	"unsafe"

	"github.com/koykov/amq"
	"github.com/koykov/openrt"
)

type Filter struct {
	amq.Base
	once sync.Once
	conf *Config

	buckets []bucket
	bp      uint64
	hsh     [256]uint64

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
	i0, i1, fp, err := f.hcalcI2FP(hkey, f.bp, 0)
	if err != nil {
		return f.mw().Set(err)
	}
	return f.hset(i0, i1, fp)
}

func (f *Filter) hset(i0, i1 uint64, fp byte) (err error) {
	b := &f.buckets[i0]
	if err = b.add(fp); err == nil {
		return f.mw().Set(nil)
	}
	b = &f.buckets[i1]
	if err = b.add(fp); err == nil {
		return f.mw().Set(nil)
	}
	i := i0
	if rand.Intn(2) == 1 {
		i = i1
	}
	for k := uint64(0); k < f.c().KicksLimit; k++ {
		j := uint64(rand.Intn(bucketsz))
		pfp := fp
		fp = f.buckets[i].fpv(j)
		_ = f.buckets[i].set(j, pfp)

		m := mask64[f.bp]
		i = (i & m) ^ (f.hsh[fp] & m)
		if err = f.buckets[i].add(fp); err == nil {
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
	return f.hset(i0, i1, fp)
}

func (f *Filter) HUnset(hkey uint64) error {
	if f.once.Do(f.init); f.err != nil {
		return f.mw().Unset(f.err)
	}
	i0, i1, fp, err := f.hcalcI2FP(hkey, f.bp, 0)
	if err != nil {
		return f.mw().Unset(err)
	}
	return f.hunset(i0, i1, fp)
}

func (f *Filter) hunset(i0, i1 uint64, fp byte) (err error) {
	b := &f.buckets[i0]
	if b.unset(fp) {
		return f.mw().Unset(nil)
	}
	b = &f.buckets[i1]
	b.unset(fp)
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
	i0, i1, fp, err := f.hcalcI2FP(hkey, f.bp, 0)
	if err != nil {
		return f.mw().Contains(false)
	}
	return f.hcontains(i0, i1, fp)
}

func (f *Filter) hcontains(i0, i1 uint64, fp byte) bool {
	b := &f.buckets[i0]
	if i := b.fpi(fp); i != -1 {
		return f.mw().Contains(true)
	}
	b = &f.buckets[i1]
	return f.mw().Contains(b.fpi(fp) != -1)
}

func (f *Filter) Size() uint64 {
	return 0 // todo implement me
}

func (f *Filter) Reset() {
	if f.once.Do(f.init); f.err != nil {
		return
	}
	openrt.MemclrUnsafe(unsafe.Pointer(&f.buckets[0]), len(f.buckets)*bucketsz)
	f.mw().Reset()
}

func (f *Filter) calcI2FP(key any, bp, i uint64) (i0 uint64, i1 uint64, fp byte, err error) {
	var hkey uint64
	if hkey, err = f.Hash(f.c().Hasher, key); err != nil {
		return
	}
	return f.hcalcI2FP(hkey, bp, i)
}

func (f *Filter) hcalcI2FP(hkey uint64, bp, i uint64) (i0 uint64, i1 uint64, fp byte, err error) {
	fp = byte(hkey%255 + 1)
	i0 = (hkey >> 32) & mask64[bp]
	m := mask64[bp]
	i1 = (i & m) ^ (f.hsh[fp] & m)
	return
}

func (f *Filter) init() {
	c := f.conf
	if c.Size == 0 {
		f.err = amq.ErrBadSize
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

	pow2 := func(n uint64) uint64 {
		n--
		n |= n >> 1
		n |= n >> 2
		n |= n >> 4
		n |= n >> 8
		n |= n >> 16
		n |= n >> 32
		n++
		return n
	}
	b := pow2(c.Size) / bucketsz
	f.bp = uint64(bits.TrailingZeros64(b))

	bc := pow2(c.Size) / bucketsz
	if bc == 0 {
		bc = 1
	}
	f.buckets = make([]bucket, bc)
	f.mw().Capacity(c.Size)

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
