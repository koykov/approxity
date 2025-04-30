package cosine

import (
	"sync"

	"github.com/koykov/byteseq"
	"github.com/koykov/pbtk/similarity"
)

type estimator[T byteseq.Q] struct {
	conf *Config[T]
	buf  []uint64
	once sync.Once

	err error
}

func NewEstimator[T byteseq.Q](conf *Config[T]) (similarity.Estimator[T], error) {
	e := &estimator[T]{conf: conf}
	if e.once.Do(e.init); e.err != nil {
		return nil, e.err
	}
	return e, nil
}

func (e *estimator[T]) Estimate(a, b T) (_ float64, err error) {
	if e.once.Do(e.init); e.err != nil {
		return 0, e.err
	}
	var mid int
	if err = e.conf.LSH.Add(a); err != nil {
		return 0, err
	}
	e.buf = e.conf.LSH.AppendHash(e.buf[:0])
	mid = len(e.buf)

	if err = e.conf.LSH.Add(b); err != nil {
		return 0, err
	}
	e.buf = e.conf.LSH.AppendHash(e.buf)

	abuf, bbuf := e.buf[:mid], e.buf[mid:]
	n := min(len(abuf), len(bbuf))
	for i := 0; i < n; i++ {
		// todo calculate...
	}
	return 0.0, nil
}

func (e *estimator[T]) Reset() {
	e.buf = e.buf[:0]
	e.conf.LSH.Reset()
}

func (e *estimator[T]) init() {
	if e.conf.LSH == nil {
		e.err = similarity.ErrNoLSH
		return
	}
}
