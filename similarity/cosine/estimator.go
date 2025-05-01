package cosine

import (
	"math"
	"sync"

	"github.com/koykov/byteseq"
	"github.com/koykov/pbtk/similarity"
)

type estimator[T byteseq.Q] struct {
	similarity.Base[T]
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
	var amag, bmag float64
	n := max(len(abuf), len(bbuf))
	_, _ = abuf[len(abuf)-1], bbuf[len(bbuf)-1]
	for i := 0; i < n; i++ {
		var av, ab float64
		if i < len(abuf) {
			av = float64(abuf[i])
		}
		if i < len(bbuf) {
			ab = float64(bbuf[i])
		}
		r += av * ab
		amag += av * av
		bmag += ab * ab
	}
	amag = math.Sqrt(amag)
	bmag = math.Sqrt(bmag)
	if amag == 0 || bmag == 0 {
		return
	}
	r /= amag * bmag
	return
}

func (e *estimator[T]) Reset() {
	e.Base.Reset()
	e.conf.LSH.Reset()
}

func (e *estimator[T]) init() {
	if e.conf.LSH == nil {
		e.err = similarity.ErrNoLSH
		return
	}
}
