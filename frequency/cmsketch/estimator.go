package cmsketch

import (
	"context"
	"io"
	"sync"

	"github.com/koykov/pbtk"
	"github.com/koykov/pbtk/frequency"
)

type estimator[T pbtk.Hashable] struct {
	pbtk.Base[T]
	conf *Config
	once sync.Once
	w, d uint64
	vec  vector

	err error
}

func NewEstimator[T pbtk.Hashable](conf *Config) (frequency.Estimator[T], error) {
	if conf == nil {
		return nil, pbtk.ErrInvalidConfig
	}
	e := &estimator[T]{conf: conf.copy()}
	if e.once.Do(e.init); e.err != nil {
		return nil, e.err
	}
	return e, nil
}

func (e *estimator[T]) Add(key T) error {
	return e.AddN(key, 1)
}

func (e *estimator[T]) AddN(key T, n uint64) error {
	if e.once.Do(e.init); e.err != nil {
		return e.err
	}
	hkey, err := e.Hash(e.conf.Hasher, key)
	if err != nil {
		return err
	}
	return e.vec.add(hkey, n)
}

func (e *estimator[T]) HAdd(hkey uint64) error {
	return e.HAddN(hkey, 1)
}

func (e *estimator[T]) HAddN(hkey uint64, n uint64) error {
	if e.once.Do(e.init); e.err != nil {
		return e.err
	}
	return e.vec.add(hkey, n)
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

func (e *estimator[T]) ReadFrom(r io.Reader) (int64, error) {
	if e.once.Do(e.init); e.err != nil {
		return 0, e.err
	}
	return e.vec.readFrom(r)
}

func (e *estimator[T]) WriteTo(w io.Writer) (int64, error) {
	if e.once.Do(e.init); e.err != nil {
		return 0, e.err
	}
	return e.vec.writeTo(w)
}

func (e *estimator[T]) Decay(ctx context.Context, factor float64) error {
	if e.once.Do(e.init); e.err != nil {
		return e.err
	}
	return e.vec.decay(ctx, factor)
}

func (e *estimator[T]) init() {
	if e.conf.Hasher == nil {
		e.err = pbtk.ErrNoHasher
		return
	}
	if e.conf.Confidence <= 0 || e.conf.Confidence >= 1 {
		e.err = frequency.ErrInvalidConfidence
		return
	}
	if e.conf.Epsilon <= 0 || e.conf.Epsilon >= 1 {
		e.err = frequency.ErrInvalidEpsilon
		return
	}
	if e.conf.MetricsWriter == nil {
		e.conf.MetricsWriter = frequency.DummyMetricsWriter{}
	}

	e.w, e.d = optimalWD(e.conf.Confidence, e.conf.Epsilon)
	if e.conf.Concurrent != nil {
		if e.conf.Compact {
			e.vec = newConcurrentVector32(e.d, e.w, e.conf.Concurrent.WriteAttemptsLimit, e.conf.Flags)
		} else {
			e.vec = newConcurrentVector64(e.d, e.w, e.conf.Concurrent.WriteAttemptsLimit, e.conf.Flags)
		}
	} else {
		if e.conf.Compact {
			e.vec = newVector32(e.d, e.w, e.conf.Flags)
		} else {
			e.vec = newVector64(e.d, e.w, e.conf.Flags)
		}
	}
}
