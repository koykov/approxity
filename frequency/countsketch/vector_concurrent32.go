package countsketch

import (
	"encoding/binary"
	"io"
	"sync/atomic"

	"github.com/koykov/pbtk"
)

const (
	dumpSignatureConcurrent32 = 0xF1B2462D8E79C8AE
	dumpVersionConcurrent32   = 1.0
)

// 32-bit version of concurrent vector implementation.
type cnvector32 struct {
	lim uint64
	buf []int32
}

func (vec *cnvector32) add(pos uint64, delta int64) error {
	for i := uint64(0); i < vec.lim+1; i++ {
		o := atomic.LoadInt32(&vec.buf[pos])
		n := o + int32(delta)
		if atomic.CompareAndSwapInt32(&vec.buf[pos], o, n) {
			return nil
		}
	}
	return pbtk.ErrWriteLimitExceed
}

func (vec *cnvector32) estimate(pos uint64) int64 {
	return int64(atomic.LoadInt32(&vec.buf[pos]))
}

func (vec *cnvector32) reset() {
	for i := uint64(0); i < uint64(len(vec.buf)); i++ {
		atomic.StoreInt32(&vec.buf[i], 0)
	}
}

func (vec *cnvector32) readFrom(r io.Reader) (n int64, err error) {
	var (
		buf [16]byte
		m   int
	)
	m, err = r.Read(buf[:])
	n += int64(m)
	if err != nil {
		return
	}

	if binary.LittleEndian.Uint64(buf[0:8]) != dumpSignatureConcurrent32 {
		err = pbtk.ErrInvalidSignature
		return
	}
	if binary.LittleEndian.Uint64(buf[8:16]) != dumpVersionConcurrent32 {
		err = pbtk.ErrVersionMismatch
		return
	}

	for i := 0; i < len(vec.buf); i++ {
		m, err = r.Read(buf[:4])
		n += int64(m)
		if err != nil {
			return
		}
		v := binary.LittleEndian.Uint32(buf[:4])
		atomic.StoreInt32(&vec.buf[i], int32(v))
	}
	return
}

func (vec *cnvector32) writeTo(w io.Writer) (n int64, err error) {
	var (
		buf [16]byte
		m   int
	)
	binary.LittleEndian.PutUint64(buf[0:8], dumpSignatureConcurrent32)
	binary.LittleEndian.PutUint64(buf[8:16], dumpVersionConcurrent32)
	m, err = w.Write(buf[:])
	n += int64(m)
	if err != nil {
		return
	}

	for i := 0; i < len(vec.buf); i++ {
		v := atomic.LoadInt32(&vec.buf[i])
		binary.LittleEndian.PutUint32(buf[0:4], uint32(v))
		m, err = w.Write(buf[:4])
		n += int64(m)
		if err != nil {
			return
		}
	}
	return
}

func newConcurrentVector32(d, w, lim uint64) *cnvector32 {
	return &cnvector32{
		lim: lim,
		buf: make([]int32, d*w),
	}
}
