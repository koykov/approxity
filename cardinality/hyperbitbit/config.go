package hyperbitbit

import (
	"github.com/koykov/approxity"
	"github.com/koykov/approxity/cardinality"
)

const defaultLgN = 5

type Config struct {
	InitialLgN    uint64
	Hasher        approxity.Hasher
	MetricsWriter cardinality.MetricsWriter
}

func NewConfig(initLgN uint64, hasher approxity.Hasher) *Config {
	return &Config{
		InitialLgN: initLgN,
		Hasher:     hasher,
	}
}

func (c *Config) copy() *Config {
	cpy := *c
	return &cpy
}
