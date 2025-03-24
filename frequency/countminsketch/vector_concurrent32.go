package countminsketch

import (
	"io"
	"sync/atomic"

	"github.com/koykov/approxity"
)

// 32-bit version of concurrent vector implementation.
type cnvector32 struct {
	basevec
	lim  uint64
	bits uint64
	buf  []uint32
}

func (vec *cnvector32) add(hkey, delta uint64) error {
	lo, hi := uint32(hkey>>32), uint32(hkey)
	for i := uint64(0); i < vec.d; i++ {
		pos := i*vec.w + uint64(lo+hi*uint32(i))%vec.w
		for j := uint64(0); j < vec.lim+1; j++ {
			o := atomic.LoadUint32(&vec.buf[pos])
			n := o + uint32(delta)
			if atomic.CompareAndSwapUint32(&vec.buf[pos], o, n) {
				return nil
			}
		}
	}
	return approxity.ErrWriteLimitExceed
}

func (vec *cnvector32) estimate(hkey uint64) (r uint64) {
	lo, hi := uint32(hkey>>32), uint32(hkey)
	for i := uint64(0); i < vec.d; i++ {
		if ce := uint64(atomic.LoadUint32(&vec.buf[vecpos(lo, hi, vec.w, i)])); r == 0 || r > ce {
			r = ce
		}
	}
	return
}

func (vec *cnvector32) reset() {
	for i := uint64(0); i < uint64(len(vec.buf)); i++ {
		atomic.StoreUint32(&vec.buf[i], 0)
	}
}

func (vec *cnvector32) readFrom(io.Reader) (n int64, err error) {
	// todo implement me
	return
}

func (vec *cnvector32) writeTo(io.Writer) (n int64, err error) {
	// todo implement me
	return
}

func newConcurrentVector32(d, w, lim uint64) *cnvector32 {
	return &cnvector32{
		basevec: basevec{d: d, w: w},
		lim:     lim,
		buf:     make([]uint32, d*w),
	}
}
