package xor

import (
	"github.com/koykov/pbtk"
	"github.com/koykov/pbtk/amq"
)

type Config struct {
	// Hasher to calculate hash sum of the items.
	// Mandatory param.
	Hasher pbtk.Hasher
	// Metrics writer handler.
	MetricsWriter amq.MetricsWriter
}

func NewConfig(hasher pbtk.Hasher) *Config {
	return &Config{Hasher: hasher}
}

func (c *Config) WithMetricsWriter(mw amq.MetricsWriter) *Config {
	c.MetricsWriter = mw
	return c
}

func (c *Config) copy() *Config {
	cpy := *c
	return &cpy
}
