package linear

import (
	"io"
	"math"
	"sync"

	"github.com/koykov/bitvector"
	"github.com/koykov/pbtk"
	"github.com/koykov/pbtk/cardinality"
)

type estimator[T pbtk.Hashable] struct {
	pbtk.Base[T]
	conf *Config
	once sync.Once
	vec  bitvector.Interface
	m    uint64

	err error
}

func NewEstimator[T pbtk.Hashable](config *Config) (cardinality.Estimator[T], error) {
	if config == nil {
		return nil, pbtk.ErrInvalidConfig
	}
	e := &estimator[T]{
		conf: config.copy(),
	}
	if e.once.Do(e.init); e.err != nil {
		return nil, e.err
	}
	return e, nil
}

func (e *estimator[T]) Add(key T) error {
	if e.once.Do(e.init); e.err != nil {
		return e.mw().Add(e.err)
	}
	hkey, err := e.Hash(e.conf.Hasher, key)
	if err != nil {
		return e.mw().Add(err)
	}
	return e.hadd(hkey)
}

func (e *estimator[T]) HAdd(hkey uint64) error {
	if e.once.Do(e.init); e.err != nil {
		return e.mw().Add(e.err)
	}
	return e.hadd(hkey)
}

func (e *estimator[T]) hadd(hkey uint64) error {
	e.vec.Set(hkey % e.m)
	return e.mw().Add(nil)
}

func (e *estimator[T]) Estimate() uint64 {
	if e.once.Do(e.init); e.err != nil {
		return e.mw().Estimate(0)
	}
	m, n := float64(e.m), float64(e.vec.Popcnt())
	return e.mw().Estimate(uint64(math.Floor(math.Abs(-m * math.Log(1-n/m)))))
}

func (e *estimator[T]) WriteTo(w io.Writer) (n int64, err error) {
	if e.once.Do(e.init); e.err != nil {
		err = e.err
		return
	}
	return e.vec.WriteTo(w)
}

func (e *estimator[T]) ReadFrom(r io.Reader) (n int64, err error) {
	if e.once.Do(e.init); e.err != nil {
		err = e.err
		return
	}
	return e.vec.ReadFrom(r)
}

func (e *estimator[T]) Reset() {
	if e.once.Do(e.init); e.err != nil {
		return
	}
	e.vec.Reset()
}

func (e *estimator[T]) init() {
	if e.conf.Hasher == nil {
		e.err = pbtk.ErrNoHasher
		return
	}
	if e.conf.MetricsWriter == nil {
		e.conf.MetricsWriter = cardinality.DummyMetricsWriter{}
	}
	if e.conf.CollisionProbability == 0 {
		e.conf.CollisionProbability = defaultCP
	}
	if e.m = optimalM(e.conf.ItemsNumber, e.conf.CollisionProbability); e.m == 0 {
		e.err = pbtk.ErrInvalidConfig
		return
	}
	if e.conf.Concurrent != nil {
		e.vec, e.err = bitvector.NewConcurrentVector(e.m, e.conf.Concurrent.WriteAttemptsLimit)
	} else {
		e.vec, e.err = bitvector.NewVector(e.m)
	}
}

func (e *estimator[T]) mw() cardinality.MetricsWriter {
	return e.conf.MetricsWriter
}
