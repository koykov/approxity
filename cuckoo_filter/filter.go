package cuckoo

import (
	"math"
	"math/bits"
	"sync"
	"unsafe"

	"github.com/koykov/amq"
	"github.com/koykov/x2bytes"
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
		return f.err
	}
	i0, i1, fp, err := f.calcI2FP(key, f.bp, 0)
	if err == nil {
		return err
	}
	if err = f.buckets[i0].add(fp); err == nil {
		return nil
	}
	if err = f.buckets[i1].add(fp); err == nil {
		return nil
	}
	for i := uint64(0); i < f.conf.KicksLimit; i++ {
		// todo implement cuckoo kicks
	}
	return nil
}

func (f *Filter) Unset(key any) error {
	if f.once.Do(f.init); f.err != nil {
		return f.err
	}
	// todo implement me
	return nil
}

func (f *Filter) Contains(key any) bool {
	if f.once.Do(f.init); f.err != nil {
		return false
	}
	// todo implement me
	return false
}

func (f *Filter) Reset() {
	if f.once.Do(f.init); f.err != nil {
		return
	}
	// todo implement me
}

func (f *Filter) calcI2FP(data any, bp, i uint64) (i0 uint64, i1 uint64, fp byte, err error) {
	const bufsz = 128
	var a [bufsz]byte
	var h struct {
		ptr      uintptr
		len, cap int
	}
	h.ptr, h.cap = uintptr(unsafe.Pointer(&a)), bufsz
	buf := *(*[]byte)(unsafe.Pointer(&h))

	if buf, err = x2bytes.ToBytes(buf, data); err != nil {
		return
	}
	hs := f.conf.Hasher.Sum64(buf)
	fp = byte(hs%255 + 1)
	i0 = (hs >> 32) & mask64[bp]
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
	if c.KicksLimit == 0 {
		c.KicksLimit = defaultKicksLimit
	}
	if c.Seed == 0 {
		c.Seed = defaultSeed
	}
	buckets := uint64(math.Ceil(float64(c.Size) / float64(8)))
	f.buckets = make([]bucket, buckets)
	f.bp = uint64(bits.TrailingZeros64(buckets))

	var buf []byte
	for i := 0; i < 256; i++ {
		buf = append(buf[:0], byte(i))
		f.hsh[i], _ = f.Hash(c.Hasher, buf, c.Seed)
	}
}
