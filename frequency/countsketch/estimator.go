package countsketch

import (
	"io"
	"slices"
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

var tsign = [2]int64{1, -1}

func NewEstimator[T pbtk.Hashable](conf *Config) (frequency.SignedEstimator[T], error) {
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
	return e.hadd(hkey, n)
}

func (e *estimator[T]) HAdd(hkey uint64) error {
	return e.HAddN(hkey, 1)
}

func (e *estimator[T]) HAddN(hkey uint64, n uint64) error {
	if e.once.Do(e.init); e.err != nil {
		return e.err
	}
	return e.hadd(hkey, n)
}

func (e *estimator[T]) hadd(hkey, n uint64) error {
	for i := uint64(0); i < e.d; i++ {
		hkeymix := e.mix(hkey, i)
		pos := hkeymix % e.w
		sign := tsign[hkeymix>>63]
		delta := sign * int64(n)
		if err := e.vec.add(i*e.w+pos, delta); err != nil {
			return err
		}
	}
	return nil
}

func (e *estimator[T]) mix(hkey, seed uint64) uint64 {
	hkey ^= seed
	hkey ^= hkey >> 33
	hkey *= 0xff51afd7ed558ccd
	hkey ^= hkey >> 33
	return hkey
}

func (e *estimator[T]) Estimate(key T) int64 {
	if e.once.Do(e.init); e.err != nil {
		return 0
	}
	hkey, err := e.Hash(e.conf.Hasher, key)
	if err != nil {
		return 0
	}
	return e.hestimate(hkey)
}

func (e *estimator[T]) HEstimate(hkey uint64) int64 {
	if e.once.Do(e.init); e.err != nil {
		return 0
	}
	return e.hestimate(hkey)
}

func (e *estimator[T]) hestimate(hkey uint64) int64 {
	var a [16]int64
	buf := a[:0]
	for i := uint64(0); i < e.d; i++ {
		hkeymix := e.mix(hkey, i)
		pos := hkeymix % e.w
		sign := tsign[hkeymix>>63]
		buf = append(buf, e.vec.estimate(i*e.w+pos)*sign)
	}
	slices.Sort(buf)
	median := buf[len(buf)/2]
	return median
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
