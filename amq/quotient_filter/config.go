package quotient

import (
	"github.com/koykov/approxity"
	"github.com/koykov/approxity/amq"
)

const (
	defaultFPP        = .01
	defaultLoadFactor = .5
)

type Config struct {
	// Number of desired items to store in the filter
	// Mandatory param.
	ItemsNumber uint64
	// False positive probability value.
	// If this param omit, defaultFPP (0.01) will use instead.
	FPP float64
	// Load factor value.
	// If this param omit, defaultLoadFactor (0.5) will use instead.
	LoadFactor float64
	// Hasher to calculate hash sum of the items.
	// Mandatory param.
	Hasher approxity.Hasher
	// Metrics writer handler.
	MetricsWriter amq.MetricsWriter
}

func NewConfig(items uint64, fpp float64, hasher approxity.Hasher) *Config {
	return &Config{
		ItemsNumber: items,
		FPP:         fpp,
		Hasher:      hasher,
	}
}

func (c *Config) WithItemsNumber(items uint64) *Config {
	c.ItemsNumber = items
	return c
}

func (c *Config) WithHasher(hasher approxity.Hasher) *Config {
	c.Hasher = hasher
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
