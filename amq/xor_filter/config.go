package xor

import (
	"github.com/koykov/approxity"
	"github.com/koykov/approxity/amq"
)

type Config struct {
	// Hasher to calculate hash sum of the items.
	// Mandatory param.
	Hasher approxity.Hasher
	// Metrics writer handler.
	MetricsWriter amq.MetricsWriter
}

func NewConfig(hasher approxity.Hasher) *Config {
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
