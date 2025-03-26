package cmsketch

import (
	"encoding/binary"
	"io"
	"unsafe"

	"github.com/koykov/approxity"
	"github.com/koykov/openrt"
)

const (
	dumpSignature64 = 0x643E037364AB8CD0
	dumpVersion64   = 1.0
)

// 64-bit version if sync vector implementation.
type syncvec64 struct {
	basevec
	buf []uint64
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

func (vec *syncvec64) readFrom(r io.Reader) (n int64, err error) {
	var (
		buf [16]byte
		m   int
	)
	m, err = r.Read(buf[:])
	n += int64(m)
	if err != nil {
		return
	}

	if binary.LittleEndian.Uint64(buf[0:8]) != dumpSignature64 {
		err = approxity.ErrInvalidSignature
		return
	}
	if binary.LittleEndian.Uint64(buf[8:16]) != dumpVersion64 {
		err = approxity.ErrVersionMismatch
		return
	}

	h := vecbufh{
		p: uintptr(unsafe.Pointer(&vec.buf[0])),
		l: len(vec.buf) * 8,
		c: len(vec.buf) * 8,
	}
	bufv := *(*[]byte)(unsafe.Pointer(&h))
	m, err = r.Read(bufv)
	n += int64(m)
	return
}

func (vec *syncvec64) writeTo(w io.Writer) (n int64, err error) {
	var (
		buf [16]byte
		m   int
	)
	binary.LittleEndian.PutUint64(buf[0:8], dumpSignature64)
	binary.LittleEndian.PutUint64(buf[8:16], dumpVersion64)
	m, err = w.Write(buf[:])
	n += int64(m)
	if err != nil {
		return
	}

	h := vecbufh{
		p: uintptr(unsafe.Pointer(&vec.buf[0])),
		l: len(vec.buf) * 8,
		c: len(vec.buf) * 8,
	}
	bufv := *(*[]byte)(unsafe.Pointer(&h))
	m, err = w.Write(bufv)
	n += int64(m)
	return
}

func newVector64(d, w uint64) *syncvec64 {
	return &syncvec64{
		basevec: basevec{d: d, w: w},
		buf:     make([]uint64, d*w),
	}
}
