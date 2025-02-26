package bloom

import (
	"github.com/koykov/amq"
	"github.com/koykov/hash"
)

type Config struct {
	// Setting up this section enables concurrent read/write operations.
	Concurrent *ConcurrentConfig
	// The size of the filter in bits.
	// Mandatory param.
	Size uint64
	// Hasher to calculate hash sum of the items.
	Hasher amq.Hasher
	// How many hash checks filter may do to reduce false positives cases.
	NumberHashFunctions uint64
	// Metrics writer handler.
	MetricsWriter amq.MetricsWriter
}

// ConcurrentConfig configures concurrent section of config.
type ConcurrentConfig struct {
	// How many write attempts may perform.
	WriteAttemptsLimit uint64
}

// NewConfig returns new config with given size and hasher.
//
// Note, size represent total size of the filter, not desired number of items, so it may be not optimal.
// Use OptimalSize function to calculate proper size and OptimalNumberHashFunctions function to calculate proper number of
// hash functions. Or use NewOptimalConfig function, it will calculate optimal params itself.
func NewConfig(size uint64, hasher amq.Hasher) *Config {
	return &Config{Size: size, Hasher: hasher}
}

// NewOptimalConfig returns new config with optimal size and number of hash functions calculated by given desired number
// of items and false positive probability.
func NewOptimalConfig(size uint64, fpp float64, hasher amq.Hasher) *Config {
	m := OptimalSize(size, fpp)
	k := OptimalNumberHashFunctions(size, m)
	return &Config{
		Size:                m,
		Hasher:              hasher,
		NumberHashFunctions: k,
	}
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

func (c *Config) WithNumberHashFunctions(number uint64) *Config {
	c.NumberHashFunctions = number
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
