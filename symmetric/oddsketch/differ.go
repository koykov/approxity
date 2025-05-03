package oddsketch

import (
	"math"
	"sync"

	"github.com/koykov/bitvector"
	"github.com/koykov/byteseq"
	"github.com/koykov/pbtk/amq"
	"github.com/koykov/pbtk/lsh"
	"github.com/koykov/pbtk/symmetric"
)

type differ[T byteseq.Q] struct {
	lsh.VectorPair[T]
	conf   *Config[T]
	m, k   uint64
	v0, v1 bitvector.Interface
	once   sync.Once

	err error
}

func NewDiffer[T byteseq.Q](conf *Config[T]) (symmetric.Differ[T], error) {
	e := &differ[T]{conf: conf.copy()}
	if e.once.Do(e.init); e.err != nil {
		return nil, e.err
	}
	return e, nil
}

func (e *differ[T]) Diff(a, b T) (r float64, err error) {
	if e.once.Do(e.init); e.err != nil {
		err = e.err
		return
	}

	abuf, bbuf, err := e.VectorizePair(e.conf.LSH, a, b)
	if len(abuf) == 0 || len(bbuf) == 0 || err != nil {
		return
	}
	for i := 0; i < len(abuf); i++ {
		e.v0.Xor(abuf[i] % e.m)
	}
	for i := 0; i < len(bbuf); i++ {
		e.v1.Xor(bbuf[i] % e.m)
	}
	var diff uint64
	if diff, err = e.v0.Difference(e.v1); err != nil || diff == 0 {
		return
	}
	m, k := float64(e.m), float64(diff)
	r = -m * math.Log(1-(2*k)/m)
	return
}

func (e *differ[T]) Reset() {
	e.VectorPair.Reset()
	e.v0.Reset()
	e.v1.Reset()
	e.conf.LSH.Reset()
}

func (e *differ[T]) init() {
	c := e.conf
	if e.conf.LSH == nil {
		e.err = symmetric.ErrNoLSH
		return
	}
	if c.ItemsNumber == 0 {
		e.err = amq.ErrNoItemsNumber
		return
	}
	if c.FPP == 0 {
		c.FPP = defaultFPP
	}
	if c.FPP < 0 || c.FPP > 1 {
		e.err = amq.ErrInvalidFPP
		return
	}

	e.m = optimalM(c.ItemsNumber, c.FPP)
	e.k = optimalK(c.ItemsNumber, e.m)

	if c.Concurrent != nil {
		if e.v0, e.err = bitvector.NewConcurrentVector(e.m, c.Concurrent.WriteAttemptsLimit); e.err != nil {
			return
		}
		e.v1, e.err = bitvector.NewConcurrentVector(e.m, c.Concurrent.WriteAttemptsLimit)
	} else {
		if e.v0, e.err = bitvector.NewVector(e.m); e.err != nil {
			return
		}
		e.v1, e.err = bitvector.NewVector(e.m)
	}
}
