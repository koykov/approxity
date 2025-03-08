package hyperloglog

import (
	"io"
	"math"

	"github.com/koykov/openrt"
)

type syncvec struct {
	a, m float64
	buf  []uint8
}

func (vec *syncvec) add(idx uint64, val uint8) error {
	o := vec.buf[idx]
	if val > o {
		vec.buf[idx] = val
	}
	return nil
}

func (vec *syncvec) estimate() (raw, nz float64) {
	buf := vec.buf
	_, _, _ = buf[len(buf)-1], pow2d1[math.MaxUint8-1], nzt[math.MaxUint8-1]
	for len(buf) > 8 {
		n0, n1, n2, n3, n4, n5, n6, n7 := buf[0], buf[1], buf[2], buf[3], buf[4], buf[5], buf[6], buf[7]
		raw += pow2d1[n0] + pow2d1[n1] + pow2d1[n2] + pow2d1[n3] + pow2d1[n4] + pow2d1[n5] + pow2d1[n6] + pow2d1[n7]
		nz += nzt[n0] + nzt[n1] + nzt[n2] + nzt[n3] + nzt[n4] + nzt[n5] + nzt[n6] + nzt[n7]
		buf = buf[8:]
	}
	for i := 0; i < len(buf); i++ {
		n := buf[i]
		raw += pow2d1[n]
		nz += nzt[n]
	}
	raw = vec.a * vec.m * vec.m / raw
	return
}

func (vec *syncvec) capacity() uint64 {
	return uint64(len(vec.buf))
}

func (vec *syncvec) reset() {
	openrt.Memclr(vec.buf)
}

func (vec *syncvec) writeTo(w io.Writer) (n int64, err error) {
	return 0, nil
}

func (vec *syncvec) readFrom(r io.Reader) (n int64, err error) {
	return 0, nil
}

func newSyncvec(a, m float64) *syncvec {
	return &syncvec{a: a, m: m, buf: make([]byte, int(m))}
}
