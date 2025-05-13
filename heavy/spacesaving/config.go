package spacesaving

import (
	"github.com/koykov/pbtk"
	"github.com/koykov/pbtk/heavy"
)

const defaultBuckets = 4

type Config struct {
	K             uint64
	Hasher        pbtk.Hasher
	Buckets       uint64
	EWMA          EWMA
	MetricsWriter heavy.MetricsWriter
}

type EWMA struct {
	Alpha float64
}

func NewConfig(k uint64, hasher pbtk.Hasher) *Config {
	return &Config{
		K:       k,
		Hasher:  hasher,
		Buckets: defaultBuckets,
	}
}

func (c *Config) WithBuckets(buckets uint64) *Config {
	c.Buckets = buckets
	return c
}

func (c *Config) WithEWMA(alpha float64) *Config {
	c.EWMA.Alpha = alpha
	return c
}

func (c *Config) WithMetricsWriter(mw heavy.MetricsWriter) *Config {
	c.MetricsWriter = mw
	return c
}

func (c *Config) copy() *Config {
	cpy := *c
	return &cpy
}
