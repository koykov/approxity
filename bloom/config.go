package bloom

import "github.com/koykov/pds"

type Config struct {
	// Setting up this section enables concurrent read/write operations.
	Concurrent *ConcurrentConfig
	// The size of the filter in bits.
	// Mandatory param.
	Size uint64
	// Hasher to calculate hash sum of the items.
	Hasher pds.Hasher
	// How many hash checks filter may do to reduce false positives cases.
	HashChecksLimit uint64
}

// ConcurrentConfig configures concurrent section of config.
type ConcurrentConfig struct {
	// How many write attempts may perform.
	WriteAttemptsLimit uint64
}

func NewConfig(size uint64, hasher pds.Hasher) *Config {
	return &Config{Size: size, Hasher: hasher}
}

func (c *Config) WithConcurrency() *Config {
	c.Concurrent = &ConcurrentConfig{}
	return c
}

func (c *Config) WithSize(size uint64) *Config {
	c.Size = size
	return c
}

func (c *Config) WithHasher(hasher pds.Hasher) *Config {
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

func (c *Config) WithHashCheckLimit(limit uint64) *Config {
	c.HashChecksLimit = limit
	return c
}

func (c *Config) copy() *Config {
	cpy := *c
	return &cpy
}
