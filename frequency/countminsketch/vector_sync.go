package countminsketch

import (
	"io"
	"unsafe"

	"github.com/koykov/openrt"
)

type syncvec[T ~uint32 | ~uint64] struct {
	d, w uint64
	buf  []T
}

func (vec *syncvec[T]) add(hkey uint64) error {
	lo, hi := uint32(hkey>>32), uint32(hkey)
	for i := uint64(0); i < vec.d; i++ {
		vec.buf[i*vec.w+uint64(lo+hi*uint32(i))%vec.w]++
	}
	return nil
}

func (vec *syncvec[T]) estimate(hkey uint64) (r uint64) {
	lo, hi := uint32(hkey>>32), uint32(hkey)
	for i := uint64(0); i < vec.d; i++ {
		if ce := uint64(vec.buf[i*vec.w+uint64(lo+hi*uint32(i))%vec.w]); r == 0 || r > ce {
			r = ce
		}
	}
	return
}

func (vec *syncvec[T]) reset() {
	sz := uint64(unsafe.Sizeof(vec.buf[0]))
	openrt.MemclrUnsafe(unsafe.Pointer(&vec.buf[0]), int(vec.w*vec.d*sz))
}

func (vec *syncvec[T]) readFrom(io.Reader) (n int64, err error) {
	// todo implement me
	return
}

func (vec *syncvec[T]) writeTo(io.Writer) (n int64, err error) {
	// todo implement me
	return
}

func newVector[T ~uint32 | ~uint64](d, w uint64) *syncvec[T] {
	return &syncvec[T]{
		d:   d,
		w:   w,
		buf: make([]T, d*w),
	}
}
