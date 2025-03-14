package xor

import (
	"github.com/koykov/approxity"
	"github.com/koykov/approxity/amq"
)

type Config struct {
	Hasher        approxity.Hasher
	MetricsWriter amq.MetricsWriter
}

func (c *Config) copy() *Config {
	cpy := *c
	return &cpy
}
