package hyperloglog

import (
	"io"
	"math"
	"math/bits"
	"sync"

	"github.com/koykov/approxity"
	"github.com/koykov/approxity/cardinality"
)

type counter[T approxity.Hashable] struct {
	approxity.Base[T]
	conf *Config
	once sync.Once
	a    float64
	m    float64
	vec  []uint8

	err error
}

func NewCounter[T approxity.Hashable](config *Config) (cardinality.Estimator[T], error) {
	if config == nil {
		return nil, approxity.ErrInvalidConfig
	}
	c := &counter[T]{
		conf: config.copy(),
	}
	if c.once.Do(c.init); c.err != nil {
		return nil, c.err
	}
	return c, nil
}

func (c *counter[T]) Add(key T) error {
	if c.once.Do(c.init); c.err != nil {
		return c.err
	}
	hkey, err := c.Hash(c.conf.Hasher, key)
	if err != nil {
		return err
	}
	return c.hadd(hkey)
}

func (c *counter[T]) HAdd(hkey uint64) error {
	if c.once.Do(c.init); c.err != nil {
		return c.err
	}
	return c.hadd(hkey)
}

func (c *counter[T]) hadd(hkey uint64) error {
	p := c.conf.Precision
	r := 64 - p
	var idx uint64
	idx = hkey >> r
	if h := hkey << p; h > 0 {
		if lz := uint64(bits.LeadingZeros64(h)) + 1; lz < r {
			r = lz
		}
	}
	c.vec[idx] = maxu8(uint8(r), c.vec[idx])
	return nil
}

func (c *counter[T]) Estimate() uint64 {
	if c.once.Do(c.init); c.err != nil || len(c.vec) == 0 {
		return 0
	}
	e, nz := c.rawEstimation()

	if e < 5*c.m {
		e = e - biasEstimation(c.conf.Precision-4, e)
	}

	h := e
	if nz < float64(uint64(1)<<c.conf.Precision) {
		h = c.linearEstimation(nz)
	}
	if h <= threshold[c.conf.Precision-4] {
		return uint64(h)
	}
	return uint64(e)
}

func (c *counter[T]) rawEstimation() (raw, nz float64) {
	_ = c.vec[len(c.vec)-1]
	for i := 0; i < len(c.vec); i++ {
		n := c.vec[i]
		raw += 1 / math.Pow(2, float64(n))
		if n > 0 {
			nz++
		}
	}
	raw = c.a * c.m * c.m / raw
	return
}

func (c *counter[T]) linearEstimation(z float64) float64 {
	return c.m * math.Log(c.m/(c.m-z))
}

func (c *counter[T]) WriteTo(w io.Writer) (n int64, err error) {
	if c.once.Do(c.init); c.err != nil {
		err = c.err
		return
	}
	// todo: implement me
	return
}

func (c *counter[T]) ReadFrom(r io.Reader) (n int64, err error) {
	if c.once.Do(c.init); c.err != nil {
		err = c.err
		return
	}
	// todo: implement me
	return
}

func (c *counter[T]) Reset() {
	if c.once.Do(c.init); c.err != nil {
		return
	}
	// todo: implement me
}

func (c *counter[T]) init() {
	if c.conf.Precision < 4 || c.conf.Precision > 18 {
		c.err = ErrInvalidPrecision
		return
	}
	if c.conf.Hasher == nil {
		c.err = approxity.ErrNoHasher
		return
	}

	m := uint64(1) << c.conf.Precision
	c.m = float64(m)
	c.vec = make([]uint8, m)

	// alpha approximation, see https://en.wikipedia.org/wiki/HyperLogLog#Practical_considerations for details
	switch m {
	case 16:
		c.a = .673
	case 32:
		c.a = .697
	case 64:
		c.a = .709
	default:
		c.a = .7213 / (1 + 1.079/c.m)
	}
}

var threshold = [15]float64{10, 20, 40, 80, 220, 400, 900, 1800, 3100, 6500, 11500, 20000, 50000, 120000, 350000}

func maxu8(a, b uint8) uint8 {
	if a > b {
		return a
	}
	return b
}
