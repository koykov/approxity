package countsketch

import (
	"github.com/koykov/approxity"
	"github.com/koykov/approxity/frequency"
)

type Config struct {
	Confidence    float64
	Epsilon       float64
	Compact       bool
	Hasher        approxity.Hasher
	Concurrent    *ConcurrentConfig
	MetricsWriter frequency.MetricsWriter
}

type ConcurrentConfig struct {
	WriteAttemptsLimit uint64
}

func NewConfig(confidence, epsilon float64, hasher approxity.Hasher) *Config {
	return &Config{
		Confidence: confidence,
		Epsilon:    epsilon,
		Hasher:     hasher,
	}
}

func (c *Config) copy() *Config {
	cpy := *c
	return &cpy
}
