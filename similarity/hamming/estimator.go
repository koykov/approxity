package hamming

import (
	"math/bits"
	"sync"

	"github.com/koykov/byteseq"
	"github.com/koykov/pbtk/lsh"
	"github.com/koykov/pbtk/similarity"
)

type estimator[T byteseq.Q] struct {
	lsh.VectorPair[T]
	conf *Config[T]
	once sync.Once

	err error
}

func NewEstimator[T byteseq.Q](conf *Config[T]) (similarity.Estimator[T], error) {
	e := &estimator[T]{conf: conf.copy()}
	if e.once.Do(e.init); e.err != nil {
		return nil, e.err
	}
	return e, nil
}

func (e *estimator[T]) Estimate(a, b T) (r float64, err error) {
	if e.once.Do(e.init); e.err != nil {
		err = e.err
		return
	}

	abuf, bbuf, err := e.VectorizePair(e.conf.LSH, a, b)
	if len(abuf) == 0 || len(bbuf) == 0 || err != nil {
		return
	}
	n := max(len(abuf), len(bbuf))
	_, _ = abuf[len(abuf)-1], bbuf[len(bbuf)-1]
	for i := 0; i < n; i++ {
		var av, bv uint64
		if i < len(abuf) {
			av = abuf[i]
		}
		if i < len(bbuf) {
			bv = bbuf[i]
		}
		r += float64(bits.OnesCount64(av ^ bv))
	}
	r = 1 - r/float64(n*64)
	return
}

func (e *estimator[T]) Reset() {
	e.VectorPair.Reset()
	e.conf.LSH.Reset()
}

func (e *estimator[T]) init() {
	if e.conf.LSH == nil {
		e.err = similarity.ErrNoLSH
		return
	}
}
