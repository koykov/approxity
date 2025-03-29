package cmsketch

// Synchronous 32/64-bit vector implementations. Generics approach is too slow in general, also there is no way
// to use atomics (in concurrent vector) together with generics.

import (
	"encoding/binary"
	"io"
	"unsafe"

	"github.com/koykov/openrt"
	"github.com/koykov/pbtk"
)

const (
	dumpSignature32 = 0x86BB26BA91E98EAD
	dumpVersion32   = 1.0
)

// 32-bit version of sync vector implementation.
type syncvec32 struct {
	basevec
	buf []uint32
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

func (vec *syncvec32) readFrom(r io.Reader) (n int64, err error) {
	var (
		buf [16]byte
		m   int
	)
	m, err = r.Read(buf[:])
	n += int64(m)
	if err != nil {
		return
	}

	if binary.LittleEndian.Uint64(buf[0:8]) != dumpSignature32 {
		err = pbtk.ErrInvalidSignature
		return
	}
	if binary.LittleEndian.Uint64(buf[8:16]) != dumpVersion32 {
		err = pbtk.ErrVersionMismatch
		return
	}

	h := vecbufh{
		p: uintptr(unsafe.Pointer(&vec.buf[0])),
		l: len(vec.buf) * 4,
		c: len(vec.buf) * 4,
	}
	bufv := *(*[]byte)(unsafe.Pointer(&h))
	m, err = r.Read(bufv)
	n += int64(m)
	return
}

func (vec *syncvec32) writeTo(w io.Writer) (n int64, err error) {
	var (
		buf [16]byte
		m   int
	)
	binary.LittleEndian.PutUint64(buf[0:8], dumpSignature32)
	binary.LittleEndian.PutUint64(buf[8:16], dumpVersion32)
	m, err = w.Write(buf[:])
	n += int64(m)
	if err != nil {
		return
	}

	h := vecbufh{
		p: uintptr(unsafe.Pointer(&vec.buf[0])),
		l: len(vec.buf) * 4,
		c: len(vec.buf) * 4,
	}
	bufv := *(*[]byte)(unsafe.Pointer(&h))
	m, err = w.Write(bufv)
	n += int64(m)
	return
}

func newVector32(d, w uint64) *syncvec32 {
	return &syncvec32{
		basevec: basevec{d: d, w: w},
		buf:     make([]uint32, d*w),
	}
}
