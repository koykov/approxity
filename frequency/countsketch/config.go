package countsketch

import (
	"github.com/koykov/approxity"
	"github.com/koykov/approxity/frequency"
)

type Config struct {
	// Confidence represent a possibility that potential error will be in range of acceptable error rate (see Epsilon).
	// E.g.: Confidence 0.99 guarantees 99% of estimations will be in range represent by Epsilon.
	// Must be in range (0..1).
	// Mandatory param.
	Confidence float64
	// Epsilon represent precision of estimation: less epsilon value makes estimation more accurate, but grows the table.
	// E.g.: Epsilon 0.01 guarantees that estimation error is less or equal 1% of all elements.
	// Must be in range (0..1).
	// Mandatory param.
	Epsilon float64
	// Hasher to calculate hash sum of the items.
	// Mandatory param.
	Hasher approxity.Hasher
	// Enable compact mode.
	// By default, uses 64 bit per counter. This param allows to use 32 bit per counter.
	Compact bool
	// Setting up this section enables concurrent read/write operations.
	Concurrent *ConcurrentConfig
	// Metrics writer handler.
	MetricsWriter frequency.MetricsWriter
}

// ConcurrentConfig configures concurrent section of config.
type ConcurrentConfig struct {
	// How many write attempts may perform.
	WriteAttemptsLimit uint64
}

func NewConfig(confidence, epsilon float64, hasher approxity.Hasher) *Config {
	return &Config{
		Confidence: confidence,
		Epsilon:    epsilon,
		Hasher:     hasher,
	}
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

func (c *Config) WithCompact() *Config {
	c.Compact = true
	return c
}

func WithMetricsWriter(conf *Config, mw frequency.MetricsWriter) *Config {
	conf.MetricsWriter = mw
	return conf
}

func (c *Config) copy() *Config {
	cpy := *c
	return &cpy
}
