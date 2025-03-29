package hyperbitbit

import (
	"github.com/koykov/pbtk"
	"github.com/koykov/pbtk/cardinality"
)

const defaultN = 1e6

type Config struct {
	// High limit of desired uniques.
	// Mandatory param.
	ItemsNumber uint64
	// Hasher to calculate hash sum of the items.
	// Mandatory param.
	Hasher pbtk.Hasher
	// Metrics writer handler.
	MetricsWriter cardinality.MetricsWriter
}

func NewConfig(itemNumber uint64, hasher pbtk.Hasher) *Config {
	return &Config{
		ItemsNumber: itemNumber,
		Hasher:      hasher,
	}
}

func (c *Config) WithMetricsWriter(mw cardinality.MetricsWriter) *Config {
	c.MetricsWriter = mw
	return c
}

func (c *Config) copy() *Config {
	cpy := *c
	return &cpy
}
