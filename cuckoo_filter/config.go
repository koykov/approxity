package cuckoo

import (
	"github.com/koykov/amq"
	"github.com/koykov/hash"
)

const defaultKicksLimit = 500

type Config struct {
	// Number of desired items to store in the filter
	// Mandatory param.
	ItemsNumber uint64
	// Hasher to calculate hash sum of the items.
	// Mandatory param.
	Hasher amq.Hasher
	// How many kicks may filter do to set the item.
	KicksLimit uint64
	// Setting up this section enables concurrent read/write operations.
	Concurrent *ConcurrentConfig
	// Metrics writer handler.
	MetricsWriter amq.MetricsWriter
}

// ConcurrentConfig configures concurrent section of config.
type ConcurrentConfig struct {
	// How many write attempts may perform.
	WriteAttemptsLimit uint64
}

func NewConfig(items uint64, hasher amq.Hasher) *Config {
	return &Config{ItemsNumber: items, Hasher: hasher}
}

func (c *Config) WithConcurrency() *Config {
	c.Concurrent = &ConcurrentConfig{}
	return c
}

func (c *Config) WithItemsNumber(items uint64) *Config {
	c.ItemsNumber = items
	return c
}

func (c *Config) WithHasher(hasher hash.Hasher[[]byte]) *Config {
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

func (c *Config) WithKicksLimit(limit uint64) *Config {
	c.KicksLimit = limit
	return c
}

func (c *Config) WithMetricsWriter(mw amq.MetricsWriter) *Config {
	c.MetricsWriter = mw
	return c
}

func (c *Config) copy() *Config {
	cpy := *c
	return &cpy
}
