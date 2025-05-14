package spacesaving

import (
	"github.com/koykov/pbtk"
	"github.com/koykov/pbtk/heavy"
)

const defaultBuckets = 4

type Config struct {
	// Number of counters.
	// Mandatory param.
	K uint64
	// Keys hasher.
	// Mandatory param.
	Hasher pbtk.Hasher
	// Number of buckets.
	// Many buckets reduces contention, but eats more memory.
	// If this param omit, defaultBuckets (4) will use instead.
	Buckets uint64
	// EWMA (Exponentially weighted moving average) settings.
	EWMA EWMA
	// Metrics writer.
	MetricsWriter heavy.MetricsWriter
}

type EWMA struct {
	// Smoothing factor.
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
