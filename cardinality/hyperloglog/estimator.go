package hyperloglog

import (
	"io"
	"math"
	"math/bits"
	"sync"

	"github.com/koykov/approxity"
	"github.com/koykov/approxity/cardinality"
	"github.com/koykov/openrt"
)

type estimator[T approxity.Hashable] struct {
	approxity.Base[T]
	conf *Config
	once sync.Once
	a    float64
	m    float64
	vec  []uint8

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
	p := e.conf.Precision
	r := 64 - p
	var idx uint64
	idx = hkey >> r
	if h := hkey << p; h > 0 {
		if lz := uint64(bits.LeadingZeros64(h)) + 1; lz < r {
			r = lz
		}
	}
	e.vec[idx] = maxu8(uint8(r), e.vec[idx])
	return nil
}

func (e *estimator[T]) Estimate() uint64 {
	if e.once.Do(e.init); e.err != nil || len(e.vec) == 0 {
		return 0
	}
	e_, nz := e.rawEstimation()

	if e_ < 5*e.m {
		e_ = e_ - biasEstimation(e.conf.Precision-4, e_)
	}

	h := e_
	if nz < float64(uint64(1)<<e.conf.Precision) {
		h = e.linearEstimation(nz)
	}
	if h <= threshold[e.conf.Precision-4] {
		return uint64(h)
	}
	return uint64(e_)
}

func (e *estimator[T]) rawEstimation() (raw, nz float64) {
	_ = e.vec[len(e.vec)-1]
	for i := 0; i < len(e.vec); i++ {
		n := e.vec[i]
		raw += 1 / math.Pow(2, float64(n))
		if n > 0 {
			nz++
		}
	}
	raw = e.a * e.m * e.m / raw
	return
}

func (e *estimator[T]) linearEstimation(z float64) float64 {
	return e.m * math.Log(e.m/(e.m-z))
}

func (e *estimator[T]) WriteTo(w io.Writer) (n int64, err error) {
	if e.once.Do(e.init); e.err != nil {
		err = e.err
		return
	}
	// todo: implement me
	return
}

func (e *estimator[T]) ReadFrom(r io.Reader) (n int64, err error) {
	if e.once.Do(e.init); e.err != nil {
		err = e.err
		return
	}
	// todo: implement me
	return
}

func (e *estimator[T]) Reset() {
	if e.once.Do(e.init); e.err != nil {
		return
	}
	openrt.Memclr(e.vec)
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
	e.vec = make([]uint8, m)

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
}

var threshold = [15]float64{10, 20, 40, 80, 220, 400, 900, 1800, 3100, 6500, 11500, 20000, 50000, 120000, 350000}

func maxu8(a, b uint8) uint8 {
	if a > b {
		return a
	}
	return b
}
