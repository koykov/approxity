package hyperbitbit

import "github.com/koykov/approxity"

const defaultLgN = 5

type Config struct {
	InitialLgN uint64
	Hasher     approxity.Hasher
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
