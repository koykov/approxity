package hyperloglog

import (
	"io"
	"math"
	"math/bits"
	"sync"

	"github.com/koykov/approxity"
	"github.com/koykov/approxity/cardinality"
)

type estimator[T approxity.Hashable] struct {
	approxity.Base[T]
	conf *Config
	once sync.Once
	a    float64
	m    float64
	vec  vector

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
	p := e.conf.Precision
	r := 64 - p
	var idx uint64
	idx = hkey >> r
	if h := hkey << p; h > 0 {
		if lz := uint64(bits.LeadingZeros64(h)) + 1; lz < r {
			r = lz
		}
	}
	return e.mw().Add(e.vec.add(idx, uint8(r)))
}

func (e *estimator[T]) Estimate() uint64 {
	if e.once.Do(e.init); e.err != nil || e.vec.capacity() == 0 {
		return 0
	}
	est, nz := e.vec.estimate()

	if est < 5*e.m {
		est = est - biasEstimation(e.conf.Precision-4, est)
	}

	h := est
	if nz < float64(uint64(1)<<e.conf.Precision) {
		h = e.linearEstimation(nz)
	}
	if h <= threshold[e.conf.Precision-4] {
		return e.mw().Estimate(uint64(h))
	}
	return e.mw().Estimate(uint64(est))
}

func (e *estimator[T]) linearEstimation(z float64) float64 {
	return e.m * math.Log(e.m/(e.m-z))
}

func (e *estimator[T]) WriteTo(w io.Writer) (n int64, err error) {
	if e.once.Do(e.init); e.err != nil {
		err = e.err
		return
	}
	return e.vec.writeTo(w)
}

func (e *estimator[T]) ReadFrom(r io.Reader) (n int64, err error) {
	if e.once.Do(e.init); e.err != nil {
		err = e.err
		return
	}
	return e.vec.readFrom(r)
}

func (e *estimator[T]) Reset() {
	if e.once.Do(e.init); e.err != nil {
		return
	}
	e.vec.reset()
}

func (e *estimator[T]) init() {
	if e.conf.Precision < 4 || e.conf.Precision > 18 {
		e.err = ErrInvalidPrecision
		return
	}
	if e.conf.Hasher == nil {
		e.err = approxity.ErrNoHasher
		return
	}
	if e.conf.MetricsWriter == nil {
		e.conf.MetricsWriter = cardinality.DummyMetricsWriter{}
	}

	m := uint64(1) << e.conf.Precision
	e.m = float64(m)

	// alpha approximation, see https://en.wikipedia.org/wiki/HyperLogLog#Practical_considerations for details
	switch m {
	case 16:
		e.a = .673
	case 32:
		e.a = .697
	case 64:
		e.a = .709
	default:
		e.a = .7213 / (1 + 1.079/e.m)
	}

	if e.conf.Concurrent != nil {
		e.vec = newCnvec(e.a, e.m, e.conf.Concurrent.WriteAttemptsLimit)
	} else {
		e.vec = newSyncvec(e.a, e.m)
	}
}

func (e *estimator[T]) mw() cardinality.MetricsWriter {
	return e.conf.MetricsWriter
}
