package bloom

import "github.com/koykov/approxity"

const defaultFPP = .01

type Config struct {
	// Number of desired items to store in the filter
	// Mandatory param.
	ItemsNumber uint64
	// False positive probability value.
	// If this param omit, defaultFPP (0.01) will use instead.
	FPP float64
	// Hasher to calculate hash sum of the items.
	// Mandatory param.
	Hasher approxity.Hasher
	// Setting up this section enables concurrent read/write operations.
	Concurrent *ConcurrentConfig
	// Metrics writer handler.
	MetricsWriter approxity.MetricsWriter
}

// ConcurrentConfig configures concurrent section of config.
type ConcurrentConfig struct {
	// How many write attempts may perform.
	WriteAttemptsLimit uint64
}

func NewConfig(items uint64, fpp float64, hasher approxity.Hasher) *Config {
	return &Config{
		ItemsNumber: items,
		FPP:         fpp,
		Hasher:      hasher,
	}
}

func (c *Config) WithConcurrency() *Config {
	c.Concurrent = &ConcurrentConfig{}
	return c
}

func (c *Config) WithItemsNumber(items uint64) *Config {
	c.ItemsNumber = items
	return c
}

func (c *Config) WithHasher(hasher approxity.Hasher) *Config {
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

func (c *Config) WithMetricsWriter(mw approxity.MetricsWriter) *Config {
	c.MetricsWriter = mw
	return c
}

func (c *Config) copy() *Config {
	cpy := *c
	return &cpy
}
