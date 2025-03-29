package loglog

import (
	"io"
	"math"
	"math/bits"
	"sync"

	"github.com/koykov/pbtk"
	"github.com/koykov/pbtk/cardinality"
)

type estimator[T pbtk.Hashable] struct {
	pbtk.Base[T]
	conf *Config
	once sync.Once

	m, a    float64
	mx, mxp uint64
	vec     vector

	err error
}

func NewEstimator[T pbtk.Hashable](conf *Config) (cardinality.Estimator[T], error) {
	if conf == nil {
		return nil, pbtk.ErrInvalidConfig
	}
	e := &estimator[T]{conf: conf}
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
	k := hkey >> e.mx
	v := uint8(bits.LeadingZeros64((hkey<<e.conf.Precision)^e.mxp)) + 1
	return e.mw().Add(e.vec.add(k, v))
}

func (e *estimator[T]) Estimate() uint64 {
	if e.once.Do(e.init); e.err != nil {
		return e.mw().Estimate(0)
	}
	raw, nz := e.vec.estimate()
	return e.mw().Estimate(uint64(e.a * e.m * (e.m - nz) / (betaEstimation(nz) + raw)))
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
		e.err = pbtk.ErrNoHasher
		return
	}
	if e.conf.MetricsWriter == nil {
		e.conf.MetricsWriter = cardinality.DummyMetricsWriter{}
	}

	m := uint64(1) << e.conf.Precision
	e.m = float64(m)
	e.mx = 64 - e.conf.Precision
	e.mxp = math.MaxUint64 >> e.mx
	e.a = .7213 / (1 + 1.079/e.m) // alpha approximation
	if e.conf.Concurrent != nil {
		e.vec = newCnvec(e.m, e.conf.Concurrent.WriteAttemptsLimit)
	} else {
		e.vec = newSyncvec(e.m)
	}
}

func (e *estimator[T]) mw() cardinality.MetricsWriter {
	return e.conf.MetricsWriter
}
