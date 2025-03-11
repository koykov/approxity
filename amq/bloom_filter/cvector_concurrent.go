package bloom

import (
	"io"
	"math"
	"sync/atomic"
)

// Concurrent counting vector implementation.
type ccnvector struct {
	buf []uint32
	lim uint64
	s   uint64
}

func (vec *ccnvector) Set(i uint64) bool {
	for j := uint64(0); j < vec.lim; j++ {
		o := atomic.LoadUint32(&vec.buf[i/2])
		v0, v1 := uint16(o>>16), uint16(o)
		if i%2 == 0 {
			v0++
		} else {
			v1++
		}
		if atomic.CompareAndSwapUint32(&vec.buf[i/2], o, uint32(v0)<<16|uint32(v1)) {
			atomic.AddUint64(&vec.s, 1)
			return true
		}
	}
	return false
}

func (vec *ccnvector) Unset(i uint64) bool {
	for j := uint64(0); j < vec.lim; j++ {
		o := atomic.LoadUint32(&vec.buf[i/2])
		v0, v1 := uint16(o>>16), uint16(o)
		if i%2 == 0 {
			v0 += math.MaxUint16
		} else {
			v1 += math.MaxUint16
		}
		if atomic.CompareAndSwapUint32(&vec.buf[i/2], o, uint32(v0)<<16|uint32(v1)) {
			atomic.AddUint64(&vec.s, math.MaxUint64)
			return true
		}
	}
	return false
}

func (vec *ccnvector) Get(i uint64) uint8 {
	c := atomic.LoadUint32(&vec.buf[i/2])
	v0, v1 := uint16(c>>16), uint16(c)
	var r bool
	if i%2 == 0 {
		r = v0 > 0
	} else {
		r = v1 > 0
	}
	if r {
		return 1
	}
	return 0
}

func (vec *ccnvector) Size() uint64 {
	return vec.s
}

func (vec *ccnvector) Capacity() uint64 {
	return uint64(len(vec.buf)) * 2
}

func (vec *ccnvector) Reset() {
	atomic.StoreUint64(&vec.s, 0)
	for i := 0; i < len(vec.buf); i++ {
		atomic.StoreUint32(&vec.buf[i], 0)
	}
}

func (vec *ccnvector) WriteTo(w io.Writer) (n int64, err error) {
	// todo implement me
	return 0, nil
}

func (vec *ccnvector) ReadFrom(r io.Reader) (n int64, err error) {
	// todo implement me
	return 0, nil
}

func newCcnvector(size, lim uint64) *ccnvector {
	return &ccnvector{
		buf: make([]uint32, size),
		lim: lim + 1,
	}
}
