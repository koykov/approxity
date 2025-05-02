package jaccard

import (
	"sync"

	"github.com/koykov/byteseq"
	"github.com/koykov/pbtk/lsh"
	"github.com/koykov/pbtk/similarity"
)

type estimator[T byteseq.Q] struct {
	lsh.VectorPair[T]
	conf   *Config[T]
	r0, r1 map[uint64]struct{}
	once   sync.Once

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
	e.rstr()
	for i := 0; i < len(abuf); i++ {
		e.r0[abuf[i]] = struct{}{}
	}
	for i := 0; i < len(bbuf); i++ {
		e.r1[bbuf[i]] = struct{}{}
	}

	var ints float64
	for h := range e.r0 {
		if _, ok := e.r1[h]; ok {
			ints++
		}
	}

	r = ints / (float64(len(e.r0)) + float64(len(e.r1)) - ints)
	return
}

func (e *estimator[T]) Reset() {
	e.VectorPair.Reset()
	e.rstr()
	e.conf.LSH.Reset()
}

func (e *estimator[T]) rstr() {
	for k := range e.r0 {
		delete(e.r0, k)
	}
	for k := range e.r1 {
		delete(e.r1, k)
	}
}

func (e *estimator[T]) init() {
	if e.conf.LSH == nil {
		e.err = similarity.ErrNoLSH
		return
	}
	e.r0 = make(map[uint64]struct{})
	e.r1 = make(map[uint64]struct{})
}
