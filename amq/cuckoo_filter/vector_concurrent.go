package cuckoo

import (
	"encoding/binary"
	"io"
	"math"
	"sync/atomic"

	"github.com/koykov/pbtk"
)

const (
	cnvecDumpSignature = 0x581fd98fe7144b7d
	cnvecDumpVersion   = 1.0
)

// Concurrent vector implementation.
type cnvec struct {
	buf []uint32
	lim uint64
	s   uint64
}

func (vec *cnvec) add(i uint64, fp byte) error {
	for k := uint64(0); k < vec.lim+1; k++ {
		for j := 0; j < bucketsz; j++ {
			if o := atomic.LoadUint32(&vec.buf[i]); o&vecmask[j] == 0 {
				n := o | uint32(fp)<<(j*8)
				if atomic.CompareAndSwapUint32(&vec.buf[i], o, n) {
					atomic.AddUint64(&vec.s, 1)
					return nil
				}
			}
		}
		return ErrFullBucket
	}
	return ErrWriteLimitReach
}

func (vec *cnvec) set(i, j uint64, fp byte) error {
	for k := uint64(0); k < vec.lim+1; k++ {
		o := atomic.LoadUint32(&vec.buf[i])
		n := o | uint32(fp)<<(j*8)
		if atomic.CompareAndSwapUint32(&vec.buf[i], o, n) {
			return nil
		}
	}
	return nil
}

func (vec *cnvec) unset(i uint64, fp byte) bool {
	for j := 0; j < bucketsz; j++ {
		if o := atomic.LoadUint32(&vec.buf[i]); o&vecmask[j] == uint32(fp)<<(j*8) {
			n := o & ^vecmask[j]
			if atomic.CompareAndSwapUint32(&vec.buf[i], o, n) {
				atomic.AddUint64(&vec.s, math.MaxUint64)
				return true
			}
		}
	}
	return false
}

func (vec *cnvec) fpv(i, j uint64) byte {
	return byte(atomic.LoadUint32(&vec.buf[i]) & vecmask[j] >> (j * 8))
}

func (vec *cnvec) fpi(i uint64, fp byte) int {
	for j := 0; j < bucketsz; j++ {
		if atomic.LoadUint32(&vec.buf[i])&vecmask[j] == uint32(fp)<<(j*8) {
			return j
		}
	}
	return -1
}

func (vec *cnvec) capacity() uint64 {
	return uint64(len(vec.buf))
}

func (vec *cnvec) size() uint64 {
	return atomic.LoadUint64(&vec.s)
}

func (vec *cnvec) reset() {
	for i := 0; i < len(vec.buf); i++ {
		atomic.StoreUint32(&vec.buf[i], 0)
	}
}

func (vec *cnvec) writeTo(w io.Writer) (n int64, err error) {
	var (
		buf [24]byte
		m   int
	)
	binary.LittleEndian.PutUint64(buf[0:8], cnvecDumpSignature)
	binary.LittleEndian.PutUint64(buf[8:16], math.Float64bits(cnvecDumpVersion))
	binary.LittleEndian.PutUint64(buf[16:24], vec.s)
	m, err = w.Write(buf[:])
	n += int64(m)
	if err != nil {
		return int64(m), err
	}

	for i := 0; i < len(vec.buf); i++ {
		v := atomic.LoadUint32(&vec.buf[i])
		var b [4]byte
		binary.LittleEndian.PutUint32(b[:], v)
		m, err = w.Write(b[:])
		n += int64(m)
		if err != nil {
			return n, err
		}
	}
	return n, err
}

func (vec *cnvec) readFrom(r io.Reader) (n int64, err error) {
	var (
		buf [24]byte
		m   int
	)
	m, err = r.Read(buf[:])
	n += int64(m)
	if err != nil {
		return n, err
	}

	sign, ver, s := binary.LittleEndian.Uint64(buf[0:8]), binary.LittleEndian.Uint64(buf[8:16]),
		binary.LittleEndian.Uint64(buf[16:24])

	if sign != cnvecDumpSignature {
		return n, pbtk.ErrInvalidSignature
	}
	if ver != math.Float64bits(cnvecDumpVersion) {
		return n, pbtk.ErrVersionMismatch
	}
	vec.s = s

	for i := 0; i < len(vec.buf); i++ {
		var b [4]byte
		m, err = r.Read(b[:])
		n += int64(m)
		if err != nil {
			if err == io.EOF {
				err = io.ErrUnexpectedEOF
			}
			return n, err
		}
		atomic.StoreUint32(&vec.buf[i], binary.LittleEndian.Uint32(b[:]))
	}
	return
}

func newCnvec(size, lim uint64) *cnvec {
	return &cnvec{
		buf: make([]uint32, size),
		lim: lim,
	}
}
