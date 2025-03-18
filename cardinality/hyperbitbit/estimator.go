package hyperbitbit

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
	conf   *Config
	once   sync.Once
	n      uint64 // lg N
	sketch [2]uint64

	err error
}

func NewEstimator[T approxity.Hashable](conf *Config) (cardinality.Estimator[T], error) {
	c := &estimator[T]{conf: conf.copy()}
	if c.once.Do(c.init); c.err != nil {
		return nil, c.err
	}
	return c, nil
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
	const d = 6
	k, z := e.klz(hkey, d)
	if z > e.n {
		e.sketch[0] |= 1 << k
	}
	if z > e.n+1 {
		e.sketch[1] |= 1 << k
	}
	if bits.OnesCount64(e.sketch[0]) > 31 {
		e.sketch[0] = e.sketch[1]
		e.sketch[1] = 0
		e.n++
	}
	return e.mw().Add(nil)
}

func (e *estimator[T]) klz(hkey uint64, d uint64) (k, z uint64) {
	m := 64 - d
	k = (hkey << 58) >> 58
	hkey >>= d
	for z = 0; hkey&0x1 == 1 && z <= m; z++ {
		hkey >>= 1
	}
	return
}

func (e *estimator[T]) hash2(hkey uint64) uint64 {
	const fib64 = 0x9e3779b97f4a7c15
	return (hkey ^ (hkey >> 32)) * fib64
}

func (e *estimator[T]) Estimate() uint64 {
	if e.once.Do(e.init); e.err != nil {
		return 0
	}
	est := float64(e.n) + 5.4 + float64(bits.OnesCount64(e.sketch[0]))/32
	est = math.Pow(2, est)
	return uint64(est)
}

func (e *estimator[T]) WriteTo(w io.Writer) (n int64, err error) {
	if e.once.Do(e.init); e.err != nil {
		err = e.err
		return
	}
	return
}

func (e *estimator[T]) ReadFrom(r io.Reader) (n int64, err error) {
	if e.once.Do(e.init); e.err != nil {
		err = e.err
		return
	}
	return
}

func (e *estimator[T]) Reset() {
	if e.once.Do(e.init); e.err != nil {
		return
	}
	e.n = e.conf.ItemsNumber
	e.sketch[0] = 0
	e.sketch[1] = 0
}

func (e *estimator[T]) init() {
	if e.conf.Hasher == nil {
		e.err = approxity.ErrNoHasher
		return
	}
	if e.conf.MetricsWriter == nil {
		e.conf.MetricsWriter = cardinality.DummyMetricsWriter{}
	}
	if e.conf.ItemsNumber == 0 {
		e.conf.ItemsNumber = defaultN
	}
	e.n = uint64(math.Log(float64(e.conf.ItemsNumber)))
}

func (e *estimator[T]) mw() cardinality.MetricsWriter {
	return e.conf.MetricsWriter
}
