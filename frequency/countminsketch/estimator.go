package countminsketch

import (
	"io"
	"sync"
	"unsafe"

	"github.com/koykov/approxity"
	"github.com/koykov/approxity/frequency"
	"github.com/koykov/openrt"
)

type estimator[T approxity.Hashable] struct {
	approxity.Base[T]
	conf *Config
	once sync.Once
	w, d uint64
	vec  []uint64

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
	return e.hadd(hkey)
}

func (e *estimator[T]) HAdd(hkey uint64) error {
	if e.once.Do(e.init); e.err != nil {
		return e.err
	}
	return e.hadd(hkey)
}

func (e *estimator[T]) hadd(hkey uint64) error {
	lo, hi := uint32(hkey>>32), uint32(hkey)
	for i := uint64(0); i < e.d; i++ {
		e.vec[uint64(lo+hi*uint32(i))%e.w]++
	}
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
	return e.hestimate(hkey)
}

func (e *estimator[T]) HEstimate(hkey uint64) uint64 {
	if e.once.Do(e.init); e.err != nil {
		return 0
	}
	return e.hestimate(hkey)
}

func (e *estimator[T]) hestimate(hkey uint64) (r uint64) {
	lo, hi := uint32(hkey>>32), uint32(hkey)
	for i := uint64(0); i < e.d; i++ {
		if ce := e.vec[uint64(lo+hi*uint32(i))%e.w]; r == 0 || r > ce {
			r = ce
		}
	}
	return
}

func (e *estimator[T]) Reset() {
	if e.once.Do(e.init); e.err != nil {
		return
	}
	openrt.MemclrUnsafe(unsafe.Pointer(&e.vec[0]), int(e.w*e.d*8))
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
	e.vec = make([]uint64, e.w*e.d)
}
