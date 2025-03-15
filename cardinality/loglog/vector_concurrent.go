package loglog

import (
	"encoding/binary"
	"io"
	"math"
	"sync/atomic"

	"github.com/koykov/approxity"
)

const (
	cnvecDumpSignature = 0xb6fcd9970391a800
	cnvecDumpVersion   = 1.0
)

type cnvec struct {
	m   float64
	lim uint64
	s   uint64
	buf []uint32
}

func (vec *cnvec) add(idx uint64, val uint8) error {
	pos, off := idx/4, idx%4
	for i := uint64(0); i < vec.lim; i++ {
		o := atomic.LoadUint32(&vec.buf[pos])
		if o8 := uint8((o >> (off * 8)) & 0xff); o8 > val {
			return nil
		}
		n := o | uint32(val)<<(off*8)
		if atomic.CompareAndSwapUint32(&vec.buf[pos], o, n) {
			atomic.AddUint64(&vec.s, 1)
			return nil
		}
	}
	return approxity.ErrWriteLimitExceed
}

func (vec *cnvec) estimate() (raw, nz float64) {
	_, _, _ = vec.buf[len(vec.buf)-1], pow2d1[math.MaxUint8-1], zt[math.MaxUint8-1]
	for i := 0; i < len(vec.buf); i++ {
		n32 := atomic.LoadUint32(&vec.buf[i])
		n0, n1, n2, n3 := n32&0xff, (n32>>8)&0xff, (n32>>16)&0xff, n32>>24
		raw += pow2d1[n0] + pow2d1[n1] + pow2d1[n2] + pow2d1[n3]
		nz += zt[n0] + zt[n1] + zt[n2] + zt[n3]
	}
	return
}

func (vec *cnvec) capacity() uint64 {
	return uint64(len(vec.buf))
}

func (vec *cnvec) size() uint64 {
	return atomic.LoadUint64(&vec.s)
}

func (vec *cnvec) reset() {
	_ = vec.buf[len(vec.buf)-1]
	for i := 0; i < len(vec.buf); i++ {
		atomic.StoreUint32(&vec.buf[i], 0)
	}
}

func (vec *cnvec) writeTo(w io.Writer) (n int64, err error) {
	var (
		buf [32]byte
		m   int
	)
	binary.LittleEndian.PutUint64(buf[0:8], cnvecDumpSignature)
	binary.LittleEndian.PutUint64(buf[8:16], math.Float64bits(cnvecDumpVersion))
	binary.LittleEndian.PutUint64(buf[16:24], math.Float64bits(vec.m))
	binary.LittleEndian.PutUint64(buf[24:32], atomic.LoadUint64(&vec.s))
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
		buf [32]byte
		m   int
	)
	m, err = r.Read(buf[:])
	n += int64(m)
	if err != nil {
		return n, err
	}

	sign, ver, m_, s := binary.LittleEndian.Uint64(buf[0:8]), binary.LittleEndian.Uint64(buf[8:16]),
		binary.LittleEndian.Uint64(buf[16:24]), binary.LittleEndian.Uint64(buf[24:32])

	if sign != cnvecDumpSignature {
		return n, approxity.ErrInvalidSignature
	}
	if ver != math.Float64bits(cnvecDumpVersion) {
		return n, approxity.ErrVersionMismatch
	}
	vec.m = math.Float64frombits(m_)
	atomic.StoreUint64(&vec.s, s)

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

func newCnvec(m float64, lim uint64) *cnvec {
	return &cnvec{m: m, lim: lim + 1, buf: make([]uint32, int(math.Ceil(m/4)))}
}
