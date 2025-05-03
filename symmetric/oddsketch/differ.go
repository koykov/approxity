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

func (d *differ[T]) Diff(a, b T) (r float64, err error) {
	if d.once.Do(d.init); d.err != nil {
		err = d.err
		return
	}

	abuf, bbuf, err := d.VectorizePair(d.conf.LSH, a, b)
	if len(abuf) == 0 || len(bbuf) == 0 || err != nil {
		return
	}
	for i := 0; i < len(abuf); i++ {
		d.v0.Xor(abuf[i] % d.m)
	}
	for i := 0; i < len(bbuf); i++ {
		d.v1.Xor(bbuf[i] % d.m)
	}
	var diff uint64
	if diff, err = d.v0.Difference(d.v1); err != nil || diff == 0 {
		return
	}
	m, k := float64(d.m), float64(diff)
	r = -m * math.Log(1-(2*k)/m)
	return
}

func (d *differ[T]) Reset() {
	d.VectorPair.Reset()
	d.v0.Reset()
	d.v1.Reset()
	d.conf.LSH.Reset()
}

func (d *differ[T]) init() {
	c := d.conf
	if d.conf.LSH == nil {
		d.err = symmetric.ErrNoLSH
		return
	}
	if c.ItemsNumber == 0 {
		d.err = amq.ErrNoItemsNumber
		return
	}
	if c.FPP == 0 {
		c.FPP = defaultFPP
	}
	if c.FPP < 0 || c.FPP > 1 {
		d.err = amq.ErrInvalidFPP
		return
	}

	d.m = optimalM(c.ItemsNumber, c.FPP)
	d.k = optimalK(c.ItemsNumber, d.m)

	if c.Concurrent != nil {
		if d.v0, d.err = bitvector.NewConcurrentVector(d.m, c.Concurrent.WriteAttemptsLimit); d.err != nil {
			return
		}
		d.v1, d.err = bitvector.NewConcurrentVector(d.m, c.Concurrent.WriteAttemptsLimit)
	} else {
		if d.v0, d.err = bitvector.NewVector(d.m); d.err != nil {
			return
		}
		d.v1, d.err = bitvector.NewVector(d.m)
	}
}
