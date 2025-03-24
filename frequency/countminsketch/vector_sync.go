package countminsketch

// Synchronous 32/64-bit vector implementations. Generics approach is too slow in general, also there is no way
// to use atomics (in concurrent vector) together with generics.

import (
	"io"
	"unsafe"

	"github.com/koykov/openrt"
)

// 32-bit version of sync vector implementation.
type syncvec32 struct {
	d, w uint64
	buf  []uint32
}

func (vec *syncvec32) add(hkey, delta uint64) error {
	lo, hi := uint32(hkey>>32), uint32(hkey)
	for i := uint64(0); i < vec.d; i++ {
		vec.buf[vecpos(lo, hi, vec.w, i)] += uint32(delta)
	}
	return nil
}

func (vec *syncvec32) estimate(hkey uint64) (r uint64) {
	lo, hi := uint32(hkey>>32), uint32(hkey)
	for i := uint64(0); i < vec.d; i++ {
		if ce := uint64(vec.buf[vecpos(lo, hi, vec.w, i)]); r == 0 || r > ce {
			r = ce
		}
	}
	return
}

func (vec *syncvec32) reset() {
	openrt.MemclrUnsafe(unsafe.Pointer(&vec.buf[0]), int(vec.w*vec.d*4))
}

func (vec *syncvec32) readFrom(io.Reader) (n int64, err error) {
	// todo implement me
	return
}

func (vec *syncvec32) writeTo(io.Writer) (n int64, err error) {
	// todo implement me
	return
}

func newVector32(d, w uint64) *syncvec32 {
	return &syncvec32{
		d:   d,
		w:   w,
		buf: make([]uint32, d*w),
	}
}

// 64-bit version if sync vector implementation.
type syncvec64 struct {
	d, w uint64
	buf  []uint64
}

func (vec *syncvec64) add(hkey, delta uint64) error {
	lo, hi := uint32(hkey>>32), uint32(hkey)
	for i := uint64(0); i < vec.d; i++ {
		vec.buf[vecpos(lo, hi, vec.w, i)] += delta
	}
	return nil
}

func (vec *syncvec64) estimate(hkey uint64) (r uint64) {
	lo, hi := uint32(hkey>>32), uint32(hkey)
	for i := uint64(0); i < vec.d; i++ {
		if ce := vec.buf[vecpos(lo, hi, vec.w, i)]; r == 0 || r > ce {
			r = ce
		}
	}
	return
}

func (vec *syncvec64) reset() {
	openrt.MemclrUnsafe(unsafe.Pointer(&vec.buf[0]), int(vec.w*vec.d*8))
}

func (vec *syncvec64) readFrom(io.Reader) (n int64, err error) {
	// todo implement me
	return
}

func (vec *syncvec64) writeTo(io.Writer) (n int64, err error) {
	// todo implement me
	return
}

func newVector64(d, w uint64) *syncvec64 {
	return &syncvec64{
		d:   d,
		w:   w,
		buf: make([]uint64, d*w),
	}
}
