package tinylfu

import (
	"github.com/koykov/pbtk"
	"github.com/koykov/pbtk/frequency"
)

type Config struct {
	Confidence    float64
	Epsilon       float64
	EWMA          EWMA
	Hasher        pbtk.Hasher
	Clock         Clock
	Concurrent    *ConcurrentConfig
	MetricsWriter frequency.MetricsWriter
}

type EWMA struct {
	Tau uint64
}

type ConcurrentConfig struct {
	WriteAttemptsLimit uint64
}

func NewConfig(confidence, epsilon float64, hasher pbtk.Hasher) *Config {
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
