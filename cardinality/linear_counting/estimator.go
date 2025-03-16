package linear

import (
	"io"
	"math"
	"sync"

	"github.com/koykov/approxity"
	"github.com/koykov/approxity/cardinality"
	"github.com/koykov/bitvector"
)

type estimator[T approxity.Hashable] struct {
	approxity.Base[T]
	conf *Config
	once sync.Once
	vec  bitvector.Interface
	m    uint64

	err error
}

func NewEstimator[T approxity.Hashable](config *Config) (cardinality.Estimator[T], error) {
	if config == nil {
		return nil, approxity.ErrInvalidConfig
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
		return e.err
	}
	hkey, err := e.Hash(e.conf.Hasher, key)
	if err != nil {
		return err
	}
	return e.hadd(hkey)
}

func (e *estimator[T]) HAdd(hkey uint64) error {
	if e.once.Do(e.init); e.err != nil {
		return e.err
	}
	return e.hadd(hkey)
}

func (e *estimator[T]) hadd(hkey uint64) error {
	e.vec.Set(hkey % e.m)
	return nil
}

func (e *estimator[T]) Estimate() uint64 {
	if e.once.Do(e.init); e.err != nil {
		return 0
	}
	m, n := float64(e.m), float64(e.vec.OnesCount())
	return uint64(math.Floor(math.Abs(-m * math.Log(1-n/m))))
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
		e.err = approxity.ErrNoHasher
		return
	}
	if e.conf.CollisionProbability == 0 {
		e.conf.CollisionProbability = defaultCP
	}
	if e.m = optimalM(e.conf.ItemsNumber, e.conf.CollisionProbability); e.m == 0 {
		e.err = approxity.ErrInvalidConfig
		return
	}
	if e.conf.Concurrent != nil {
		e.vec, e.err = bitvector.NewConcurrentVector(e.m, e.conf.Concurrent.WriteAttemptsLimit)
	} else {
		e.vec, e.err = bitvector.NewVector(e.m)
	}
}
