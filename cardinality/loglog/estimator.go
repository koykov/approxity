package loglog

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

	m, a    float64
	mx, mxp uint64
	vec     vector

	err error
}

func NewEstimator[T approxity.Hashable](conf *Config) (cardinality.Estimator[T], error) {
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
	k := hkey >> e.mx
	v := uint8(bits.LeadingZeros64((hkey<<e.conf.Precision)^e.mxp)) + 1
	return e.vec.add(k, v)
}

func (e *estimator[T]) Estimate() uint64 {
	if e.once.Do(e.init); e.err != nil {
		return 0
	}
	raw, nz := e.vec.estimate()
	return uint64(e.a * e.m * (e.m - nz) / (betaEstimation(nz) + raw))
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
