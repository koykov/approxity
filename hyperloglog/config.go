package hyperloglog

import "github.com/koykov/amq"

type Config struct {
	// Must be in range [4..18].
	Precision      uint64
	BiasCorrection bool
	Hasher         amq.Hasher
}

func (c *Config) copy() *Config {
	cpy := *c
	return &cpy
}
