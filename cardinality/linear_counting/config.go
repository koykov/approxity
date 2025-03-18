package linear

import (
	"github.com/koykov/approxity"
	"github.com/koykov/approxity/cardinality"
)

const defaultCP = .01

type Config struct {
	// High limit of desired uniques.
	// Mandatory param.
	ItemsNumber uint64
	// Collision probability in range (0..1).
	// If this param omit, defaultCP (0.01) will use instead.
	CollisionProbability float64
	// Hasher to calculate hash sum of the items.
	// Mandatory param.
	Hasher approxity.Hasher
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

func NewConfig(itemsNumber uint64, hasher approxity.Hasher) *Config {
	return &Config{
		ItemsNumber: itemsNumber,
		Hasher:      hasher,
	}
}

func (c *Config) WithCollisionProbability(p float64) *Config {
	c.CollisionProbability = p
	return c
}

func (c *Config) WithConcurrency() *Config {
	c.Concurrent = &ConcurrentConfig{}
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
