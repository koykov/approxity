package hyperloglog

import (
	"github.com/koykov/pbtk"
	"github.com/koykov/pbtk/cardinality"
)

type Config struct {
	// Must be in range [4..18].
	// Mandatory param.
	Precision uint64
	// Hasher to calculate hash sum of the items.
	// Mandatory param.
	Hasher pbtk.Hasher
	// Setting up this section enables concurrent read/write operations.
	Concurrent *ConcurrentConfig
	// Metrics writer handler.
	MetricsWriter cardinality.MetricsWriter
}

// ConcurrentConfig configures concurrent section of config.
type ConcurrentConfig struct {
	// How many write attempts may perform.
	WriteAttemptsLimit uint64
}

func NewConfig(precision uint64, hasher pbtk.Hasher) *Config {
	return &Config{
		Precision: precision,
		Hasher:    hasher,
	}
}

func (c *Config) WithConcurrency() *Config {
	c.Concurrent = &ConcurrentConfig{}
	return c
}

func (c *Config) WithPrecision(precision uint64) *Config {
	c.Precision = precision
	return c
}

func (c *Config) WithHasher(hasher pbtk.Hasher) *Config {
	c.Hasher = hasher
	return c
}

func (c *Config) WithWriteAttemptsLimit(limit uint64) *Config {
	if c.Concurrent == nil {
		c.Concurrent = &ConcurrentConfig{}
	}
	c.Concurrent.WriteAttemptsLimit = limit
	return c
}

func (c *Config) WithMetricsWriter(mw cardinality.MetricsWriter) *Config {
	c.MetricsWriter = mw
	return c
}

func (c *Config) copy() *Config {
	cpy := *c
	return &cpy
}
