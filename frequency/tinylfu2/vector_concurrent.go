package tinylfu

import (
	"io"
	"sync/atomic"

	"github.com/koykov/pbtk"
)

type cnvec struct {
	basevec
	lim uint64
}

func (vec *cnvec) set(pos, n uint64, dtime uint32) error {
	for i := uint64(0); i < vec.lim+1; i++ {
		val := atomic.LoadUint64(&vec.buf[pos])
		newVal := vec.recalc(val, n, dtime)
		if atomic.CompareAndSwapUint64(&vec.buf[pos], val, newVal) {
			return nil
		}
	}
	return pbtk.ErrWriteLimitExceed
}

func (vec *cnvec) get(pos uint64, stime, now uint32) uint32 {
	val := atomic.LoadUint64(&vec.buf[pos])
	return vec.estimate(val, stime, now)
}

func (vec *cnvec) reset() {
	for i := 0; i < len(vec.buf); i++ {
		atomic.StoreUint64(&vec.buf[i], 0)
	}
}

func (vec *cnvec) readFrom(r io.Reader) (int64, error) {
	// todo implement me
	return 0, nil
}

func (vec *cnvec) writeTo(w io.Writer) (int64, error) {
	// todo implement me
	return 0, nil
}

func newConcurrentVector(sz, lim uint64, ewma *EWMA) vector {
	return &cnvec{
		basevec: basevec{
			buf:      make([]uint64, sz),
			dtimeMin: ewma.MinDeltaTime,
			tau:      ewma.Tau,
		},
		lim: lim,
	}
}
