package tinylfu

import (
	"io"
	"unsafe"

	"github.com/koykov/openrt"
)

type syncvec struct {
	*basevec
}

func (vec *syncvec) set(pos, n uint64, dtime uint32) error {
	val := vec.buf[pos]
	val = vec.recalc(val, n, dtime)
	vec.buf[pos] = val
	return nil
}

func (vec *syncvec) get(pos uint64, stime, now uint32) uint32 {
	val := vec.buf[pos]
	return vec.estimate(val, stime, now)
}

func (vec *syncvec) reset() {
	openrt.MemclrUnsafe(unsafe.Pointer(&vec.buf[0]), len(vec.buf)*8)
}

func (vec *syncvec) readFrom(r io.Reader) (int64, error) {
	// todo implement me
	return 0, nil
}

func (vec *syncvec) writeTo(w io.Writer) (int64, error) {
	// todo implement me
	return 0, nil
}

func newVector(sz uint64, ewma *EWMA) vector {
	vec := &syncvec{
		basevec: &basevec{
			buf:      make([]uint64, sz),
			dtimeMin: ewma.MinDeltaTime,
			tau:      ewma.Tau,
			exptabsz: ewma.ExpTableSize,
		},
	}
	vec.basevec.init()
	return vec
}
