package cuckoo

import (
	"math/bits"
	"math/rand"
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
	if err != nil {
		return err
	}
	b := &f.buckets[i0]
	if err = b.add(fp); err == nil {
		return nil
	}
	b = &f.buckets[i1]
	if err = b.add(fp); err == nil {
		return nil
	}
	i := i0
	if rand.Intn(2) == 1 {
		i = i1
	}
	for k := uint64(0); k < f.conf.KicksLimit; k++ {
		j := uint64(rand.Intn(bucketsz))
		pfp := fp
		fp = f.buckets[i].fpv(j)
		_ = f.buckets[i].set(j, pfp)

		m := mask64[f.bp]
		i = (i & m) ^ (f.hsh[fp] & m)
		if err = f.buckets[i].add(fp); err == nil {
			return nil
		}
	}
	return ErrFullFilter
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
	i0, i1, fp, err := f.calcI2FP(key, f.bp, 0)
	if err != nil {
		return false
	}
	b := &f.buckets[i0]
	if i := b.fpi(fp); i != -1 {
		return true
	}
	b = &f.buckets[i1]
	return b.fpi(fp) != -1
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

	var buf []byte
	for i := 0; i < 256; i++ {
		buf = append(buf[:0], byte(i))
		f.hsh[i], _ = f.Hash(c.Hasher, buf, c.Seed)
	}
}
