package cuckoo

import (
	"github.com/koykov/amq"
	"github.com/koykov/hash"
)

const defaultKicksLimit = 500

type Config struct {
	Concurrent    *ConcurrentConfig
	Size          uint64
	KicksLimit    uint64
	Hasher        amq.Hasher
	MetricsWriter amq.MetricsWriter
}

// ConcurrentConfig configures concurrent section of config.
type ConcurrentConfig struct {
	// How many write attempts may perform.
	WriteAttemptsLimit uint64
}

func (c *Config) WithConcurrency() *Config {
	c.Concurrent = &ConcurrentConfig{}
	return c
}

func (c *Config) WithSize(size uint64) *Config {
	c.Size = size
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
