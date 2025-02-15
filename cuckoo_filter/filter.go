package cuckoo

import (
	"math"
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
		conf: conf.copy(),
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
	buckets := uint64(math.Ceil(float64(c.Size) / float64(c.BucketSize)))
	f.buckets = make([]bucket, buckets)
	f.buf = make([]byte, c.Size*c.FingerprintSize)
}

func CalcFingerprintSize(size uint64, fp float64) (sz uint64) {
	if sz = uint64(math.Ceil(math.Log(2*float64(size)/fp))) / 8; sz == 0 {
		sz = 1
	}
	return
}
