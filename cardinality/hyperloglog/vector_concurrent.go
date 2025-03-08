package hyperloglog

import (
	"io"
	"math"
	"sync/atomic"

	"github.com/koykov/approxity/cardinality"
)

type cnvec struct {
	a, m float64
	lim  uint64
	buf  []uint32
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
			return nil
		}
	}
	return cardinality.ErrWriteAttemptsLimitExceeded
}

func (vec *cnvec) estimate() (raw, nz float64) {
	// _, _, _ = vec.buf[len(vec.buf)-1], pow2d1[math.MaxUint8-1], nzt[math.MaxUint8-1]
	for i := 0; i < len(vec.buf); i++ {
		n32 := atomic.LoadUint32(&vec.buf[i])
		n0, n1, n2, n3 := n32&0xff, (n32>>8)&0xff, (n32>>16)&0xff, n32>>24
		raw += pow2d1[n0] + pow2d1[n1] + pow2d1[n2] + pow2d1[n3]
		nz += nzt[n0] + nzt[n1] + nzt[n2] + nzt[n3]
	}
	raw = vec.a * vec.m * vec.m / raw
	return
}

func (vec *cnvec) capacity() uint64 {
	return uint64(len(vec.buf))
}

func (vec *cnvec) reset() {
	_ = vec.buf[len(vec.buf)-1]
	for i := 0; i < len(vec.buf); i++ {
		atomic.StoreUint32(&vec.buf[i], 0)
	}
}

func (vec *cnvec) writeTo(w io.Writer) (n int64, err error) {
	return 0, nil
}

func (vec *cnvec) readFrom(r io.Reader) (n int64, err error) {
	return 0, nil
}

func newCnvec(a, m float64, lim uint64) *cnvec {
	return &cnvec{a: a, m: m, lim: lim + 1, buf: make([]uint32, int(math.Ceil(m/4)))}
}
