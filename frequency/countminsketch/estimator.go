package countminsketch

import (
	"io"
	"sync"

	"github.com/koykov/approxity"
	"github.com/koykov/approxity/frequency"
)

type estimator[T approxity.Hashable] struct {
	approxity.Base[T]
	conf *Config
	once sync.Once
	w, d uint64
	vec  vector

	err error
}

func NewEstimator[T approxity.Hashable](conf *Config) (frequency.Estimator[T], error) {
	if conf == nil {
		return nil, approxity.ErrInvalidConfig
	}
	e := &estimator[T]{conf: conf}
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
	return e.vec.add(hkey, 1)
}

func (e *estimator[T]) HAdd(hkey uint64) error {
	if e.once.Do(e.init); e.err != nil {
		return e.err
	}
	return e.vec.add(hkey, 1)
}

func (e *estimator[T]) Estimate(key T) uint64 {
	if e.once.Do(e.init); e.err != nil {
		return 0
	}
	hkey, err := e.Hash(e.conf.Hasher, key)
	if err != nil {
		return 0
	}
	return e.vec.estimate(hkey)
}

func (e *estimator[T]) HEstimate(hkey uint64) uint64 {
	if e.once.Do(e.init); e.err != nil {
		return 0
	}
	return e.vec.estimate(hkey)
}

func (e *estimator[T]) Reset() {
	if e.once.Do(e.init); e.err != nil {
		return
	}
	e.vec.reset()
}

func (e *estimator[T]) ReadFrom(r io.Reader) (n int64, err error) {
	if e.once.Do(e.init); e.err != nil {
		err = e.err
		return
	}
	return
}

func (e *estimator[T]) WriteTo(w io.Writer) (n int64, err error) {
	if e.once.Do(e.init); e.err != nil {
		err = e.err
		return
	}
	return
}

func (e *estimator[T]) init() {
	if e.conf.Hasher == nil {
		e.err = approxity.ErrNoHasher
		return
	}
	if e.conf.Confidence == 0 || e.conf.Confidence > 1 {
		e.err = ErrInvalidConfidence
		return
	}
	if e.conf.Epsilon == 0 || e.conf.Epsilon > 1 {
		e.err = ErrInvalidEpsilon
		return
	}
	if e.conf.MetricsWriter == nil {
		e.conf.MetricsWriter = frequency.DummyMetricsWriter{}
	}

	e.w, e.d = optimalWD(e.conf.Confidence, e.conf.Epsilon)
	if e.conf.Concurrent != nil {
		if e.conf.Compact {
			e.vec = newConcurrentVector32(e.d, e.w, e.conf.Concurrent.WriteAttemptsLimit)
		} else {
			e.vec = newConcurrentVector64(e.d, e.w, e.conf.Concurrent.WriteAttemptsLimit)
		}
	} else {
		if e.conf.Compact {
			e.vec = newVector32(e.d, e.w)
		} else {
			e.vec = newVector64(e.d, e.w)
		}
	}
}
