package hyperloglog

import (
	"github.com/koykov/approxity"
	"github.com/koykov/approxity/cardinality"
)

type Config struct {
	// Must be in range [4..18].
	Precision     uint64
	Hasher        approxity.Hasher
	MetricsWriter cardinality.MetricsWriter
}

func (c *Config) copy() *Config {
	cpy := *c
	return &cpy
}
