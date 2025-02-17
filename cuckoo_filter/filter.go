package cuckoo

import (
	"math"
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
	buf     []byte
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
	// todo implement me
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
	if c.BucketSize == 0 {
		c.BucketSize = defaultBucketSize
	}
	if c.FingerprintSize == 0 {
		c.FingerprintSize = CalcFingerprintSize(c.Size, defaultFalsePositiveRate)
	}
	if c.KicksLimit == 0 {
		c.KicksLimit = defaultKicksLimit
	}
	if c.Seed == 0 {
		c.Seed = defaultSeed
	}
	buckets := uint64(math.Ceil(float64(c.Size) / float64(c.BucketSize)))
	f.buckets = make([]bucket, buckets)
	f.buf = make([]byte, c.Size*c.FingerprintSize)

	var buf []byte
	for i := 0; i < 256; i++ {
		buf = append(buf[:0], byte(i))
		f.hsh[i], _ = f.Hash(c.Hasher, buf, c.Seed)
	}
}

func CalcFingerprintSize(size uint64, fp float64) (sz uint64) {
	if sz = uint64(math.Ceil(math.Log(2*float64(size)/fp))) / 8; sz == 0 {
		sz = 1
	}
	return
}
