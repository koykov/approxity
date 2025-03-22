package countminsketch

import (
	"github.com/koykov/approxity"
	"github.com/koykov/approxity/frequency"
)

type Config struct {
	Epsilon       float64
	Confidence    float64
	Hasher        approxity.Hasher
	MetricsWriter frequency.MetricsWriter
}

func NewConfig(confidence, epsilon float64, hasher approxity.Hasher) *Config {
	return &Config{
		Confidence: confidence,
		Epsilon:    epsilon,
		Hasher:     hasher,
	}
}

func WithMetricsWriter(conf *Config, mw frequency.MetricsWriter) *Config {
	conf.MetricsWriter = mw
	return conf
}

func (c *Config) copy() *Config {
	cpy := *c
	return &cpy
}
