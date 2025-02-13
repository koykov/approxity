package bloom

import (
	"strconv"
	"sync"
	"unsafe"

	"github.com/koykov/bitvector"
	"github.com/koykov/x2bytes"
)

// Filter represents Bloom filter.
// By default, filter doesn't support concurrent read/write operations - you must set up the filter before reading.
// Concurrent reading allowed afterward.
// If you want to use concurrent read/write operations, fill up Concurrent section in Config object.
type Filter struct {
	once sync.Once
	conf *Config
	vec  bitvector.Interface

	err error
}

// NewFilter creates new Bloom filter.
func NewFilter(config *Config) (*Filter, error) {
	if config == nil {
		return nil, ErrBadConfig
	}
	f := &Filter{
		conf: config.copy(),
	}
	f.once.Do(f.init)
	return f, f.err
}

// Set adds new key to the filter.
func (f *Filter) Set(key any) error {
	if f.once.Do(f.init); f.err != nil {
		return f.err
	}
	for i := uint64(0); i < f.conf.HashChecksLimit+1; i++ {
		h, err := f.hash(key, i)
		if err != nil {
			return err
		}
		f.vec.Set(h % f.conf.Size)
	}
	return nil
}

// Unset removes key from the filter.
func (f *Filter) Unset(key any) error {
	if f.once.Do(f.init); f.err != nil {
		return f.err
	}
	for i := uint64(0); i < f.conf.HashChecksLimit+1; i++ {
		h, err := f.hash(key, i)
		if err != nil {
			return err
		}
		f.vec.Unset(h % f.conf.Size)
	}
	return nil
}

// Contains checks if key is in the filter.
func (f *Filter) Contains(key any) bool {
	if f.once.Do(f.init); f.err != nil {
		return false
	}
	for i := uint64(0); i < f.conf.HashChecksLimit+1; i++ {
		h, err := f.hash(key, i)
		if err != nil {
			return false
		}
		if f.vec.Get(h%f.conf.Size) == 0 {
			return false
		}
	}
	return true
}

func (f *Filter) Reset() {
	if f.once.Do(f.init); f.err != nil {
		return
	}
	f.vec.Reset()
}

func (f *Filter) hash(data any, seed uint64) (_ uint64, err error) {
	const bufsz = 128
	var a [bufsz]byte
	var h struct {
		ptr      uintptr
		len, cap int
	}
	h.ptr, h.cap = uintptr(unsafe.Pointer(&a)), bufsz
	buf := *(*[]byte)(unsafe.Pointer(&h))

	if buf, err = x2bytes.ToBytes(buf, data); err != nil {
		return 0, err
	}
	buf = strconv.AppendUint(buf, seed, 10)
	return f.conf.Hasher.Sum64(buf), nil
}

func (f *Filter) init() {
	if f.conf.Hasher == nil {
		f.err = ErrNoHasher
		return
	}
	if f.conf.Concurrent != nil {
		f.vec, f.err = bitvector.NewConcurrentVector(f.conf.Size, f.conf.Concurrent.WriteAttemptsLimit)
	} else {
		f.vec, f.err = bitvector.NewVector(f.conf.Size)
	}
}
