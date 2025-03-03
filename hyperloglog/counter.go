package hyperloglog

import (
	"io"
	"math"
	"math/bits"
	"sync"

	"github.com/koykov/amq"
)

type counter struct {
	amq.Base
	conf *Config
	once sync.Once
	a    float64
	// b    uint64
	m    uint64
	mpw2 uint64
	vec  []uint8

	err error
}

func NewCounter(config *Config) (amq.Counter, error) {
	if config == nil {
		return nil, amq.ErrInvalidConfig
	}
	c := &counter{
		conf: config.copy(),
	}
	if c.once.Do(c.init); c.err != nil {
		return nil, c.err
	}
	return c, nil
}

func (c *counter) Add(key any) error {
	if c.once.Do(c.init); c.err != nil {
		return c.err
	}
	hkey, err := c.Hash(c.conf.Hasher, key)
	if err != nil {
		return err
	}
	return c.hadd(hkey)
}

func (c *counter) HAdd(hkey uint64) error {
	if c.once.Do(c.init); c.err != nil {
		return c.err
	}
	return c.hadd(hkey)
}

func (c *counter) hadd(hkey uint64) error {
	p := c.conf.Precision
	i := (hkey >> (64 - p)) & ((1 << p) - 1)
	w := hkey<<p | 1<<(p-1)
	lbp1 := uint8(bits.LeadingZeros64(w)) + 1
	if mx := c.vec[i]; lbp1 > mx {
		c.vec[i] = lbp1
	}
	return nil
}

func (c *counter) Count() uint64 {
	if c.once.Do(c.init); c.err != nil || len(c.vec) == 0 {
		return 0
	}
	var s, z float64
	_ = c.vec[len(c.vec)-1]
	for i := 0; i < len(c.vec); i++ {
		n := c.vec[i]
		s += 1 / math.Pow(2, float64(n))
		if n == 0 {
			z++
		}
	}

	m := float64(c.m)
	r := math.Ceil(c.a * m * (m - z) / (s + betafn(c.conf.Precision, z)))

	if c.conf.BiasCorrection {
		r = biasfn(c.conf.Precision, r)
	}

	return uint64(r)
}

func (c *counter) WriteTo(w io.Writer) (n int64, err error) {
	if c.once.Do(c.init); c.err != nil {
		err = c.err
		return
	}
	// todo: implement me
	return
}

func (c *counter) ReadFrom(r io.Reader) (n int64, err error) {
	if c.once.Do(c.init); c.err != nil {
		err = c.err
		return
	}
	// todo: implement me
	return
}

func (c *counter) Reset() {
	if c.once.Do(c.init); c.err != nil {
		return
	}
	// todo: implement me
}

func (c *counter) init() {
	if c.conf.Precision < 4 || c.conf.Precision > 18 {
		c.err = ErrInvalidPrecision
		return
	}
	if c.conf.Hasher == nil {
		c.err = amq.ErrNoHasher
		return
	}

	c.m = 1 << c.conf.Precision
	c.mpw2 = c.m * c.m

	// alpha approximation, see https://en.wikipedia.org/wiki/HyperLogLog#Practical_considerations for details
	switch c.m {
	case 16:
		c.a = .673
	case 32:
		c.a = .697
	case 64:
		c.a = .709
	default:
		c.a = .7213 / (1 + 1.079/float64(c.m))
	}
}
