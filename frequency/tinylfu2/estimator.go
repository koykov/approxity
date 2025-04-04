package tinylfu

import (
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
	vec  []uint64

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
	_ = hkey
	// todo implement me
	return nil
}

func (e *estimator[T]) HAdd(hkey uint64) error {
	return e.HAddN(hkey, 1)
}

func (e *estimator[T]) HAddN(hkey uint64, n uint64) error {
	if e.once.Do(e.init); e.err != nil {
		return e.err
	}
	// todo implement me
	return nil
}

func (e *estimator[T]) Estimate(key T) uint64 {
	if e.once.Do(e.init); e.err != nil {
		return 0
	}
	hkey, err := e.Hash(e.conf.Hasher, key)
	if err != nil {
		return 0
	}
	_ = hkey
	// todo implement me
	return 0
}

func (e *estimator[T]) HEstimate(hkey uint64) uint64 {
	if e.once.Do(e.init); e.err != nil {
		return 0
	}
	// todo implement me
	return 0
}

func (e *estimator[T]) Reset() {
	if e.once.Do(e.init); e.err != nil {
		return
	}
	// todo implement me
}

func (e *estimator[T]) ReadFrom(r io.Reader) (int64, error) {
	if e.once.Do(e.init); e.err != nil {
		return 0, e.err
	}
	// todo implement me
	return 0, nil
}

func (e *estimator[T]) WriteTo(w io.Writer) (int64, error) {
	if e.once.Do(e.init); e.err != nil {
		return 0, e.err
	}
	// todo implement me
	return 0, nil
}

func (e *estimator[T]) init() {
	// todo implement me
}
