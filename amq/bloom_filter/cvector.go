package bloom

import (
	"io"
	"math"
	"unsafe"

	"github.com/koykov/openrt"
)

// Synchronous counting vector implementation.
type cvector struct {
	buf []uint32
	s   uint64
}

func (vec *cvector) Set(i uint64) bool {
	c := vec.buf[i/2]
	v0, v1 := uint16(c>>16), uint16(c)
	if i%2 == 0 {
		v0++
	} else {
		v1++
	}
	vec.buf[i/2] = uint32(v0)<<16 | uint32(v1)
	vec.s++
	return true
}

func (vec *cvector) Unset(i uint64) bool {
	c := vec.buf[i/2]
	v0, v1 := uint16(c>>16), uint16(c)
	if i%2 == 0 {
		v0 += math.MaxUint16
	} else {
		v1 += math.MaxUint16
	}
	vec.buf[i/2] = uint32(v0)<<16 | uint32(v1)
	vec.s += math.MaxUint16
	return true
}

func (vec *cvector) Get(i uint64) uint8 {
	c := vec.buf[i/2]
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

func (vec *cvector) Size() uint64 {
	return vec.s
}

func (vec *cvector) Capacity() uint64 {
	return uint64(len(vec.buf)) * 2
}

func (vec *cvector) Reset() {
	openrt.MemclrUnsafe(unsafe.Pointer(&vec.buf[0]), len(vec.buf)*4)
	vec.s = 0
}

func (vec *cvector) WriteTo(w io.Writer) (n int64, err error) {
	// todo implement me
	return 0, nil
}

func (vec *cvector) ReadFrom(r io.Reader) (n int64, err error) {
	// todo implement me
	return 0, nil
}

func newCvector(size uint64) *cvector {
	return &cvector{buf: make([]uint32, size)}
}
