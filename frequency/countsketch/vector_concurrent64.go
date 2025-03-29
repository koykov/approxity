package countsketch

import (
	"encoding/binary"
	"io"
	"sync/atomic"

	"github.com/koykov/approxity"
)

const (
	dumpSignatureConcurrent64 = 0x344E143615C13A0A
	dumpVersionConcurrent64   = 1.0
)

// 64-bit version of concurrent vector implementation.
type cnvector64 struct {
	lim uint64
	buf []int64
}

func (vec *cnvector64) add(pos uint64, delta int64) error {
	for i := uint64(0); i < vec.lim+1; i++ {
		o := atomic.LoadInt64(&vec.buf[pos])
		n := o + delta
		if atomic.CompareAndSwapInt64(&vec.buf[pos], o, n) {
			return nil
		}
	}
	return approxity.ErrWriteLimitExceed
}

func (vec *cnvector64) estimate(pos uint64) int64 {
	return atomic.LoadInt64(&vec.buf[pos])
}

func (vec *cnvector64) reset() {
	for i := uint64(0); i < uint64(len(vec.buf)); i++ {
		atomic.StoreInt64(&vec.buf[i], 0)
	}
}

func (vec *cnvector64) readFrom(r io.Reader) (n int64, err error) {
	var (
		buf [16]byte
		m   int
	)
	m, err = r.Read(buf[:])
	n += int64(m)
	if err != nil {
		return
	}

	if binary.LittleEndian.Uint64(buf[0:8]) != dumpSignatureConcurrent64 {
		err = approxity.ErrInvalidSignature
		return
	}
	if binary.LittleEndian.Uint64(buf[8:16]) != dumpVersionConcurrent64 {
		err = approxity.ErrVersionMismatch
		return
	}

	for i := 0; i < len(vec.buf); i++ {
		m, err = r.Read(buf[:8])
		n += int64(m)
		if err != nil {
			return
		}
		v := binary.LittleEndian.Uint64(buf[:8])
		atomic.StoreInt64(&vec.buf[i], int64(v))
	}
	return
}

func (vec *cnvector64) writeTo(w io.Writer) (n int64, err error) {
	var (
		buf [16]byte
		m   int
	)
	binary.LittleEndian.PutUint64(buf[0:8], dumpSignatureConcurrent64)
	binary.LittleEndian.PutUint64(buf[8:16], dumpVersionConcurrent64)
	m, err = w.Write(buf[:])
	n += int64(m)
	if err != nil {
		return
	}

	for i := 0; i < len(vec.buf); i++ {
		v := atomic.LoadInt64(&vec.buf[i])
		binary.LittleEndian.PutUint64(buf[0:8], uint64(v))
		m, err = w.Write(buf[:8])
		n += int64(m)
		if err != nil {
			return
		}
	}
	return
}

func newConcurrentVector64(d, w, lim uint64) *cnvector64 {
	return &cnvector64{
		lim: lim,
		buf: make([]int64, d*w),
	}
}
