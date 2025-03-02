package quotient

import "github.com/koykov/amq"

const (
	defaultFPP        = .01
	defaultLoadFactor = .5
)

type Config struct {
	ItemsNumber   uint64
	FPP           float64
	LoadFactor    float64
	Hasher        amq.Hasher
	MetricsWriter amq.MetricsWriter
}

func NewConfig(items uint64, fpp float64, hasher amq.Hasher) *Config {
	return &Config{
		ItemsNumber: items,
		FPP:         fpp,
		Hasher:      hasher,
	}
}

func (c *Config) copy() *Config {
	cpy := *c
	return &cpy
}
