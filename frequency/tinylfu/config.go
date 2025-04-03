package tinylfu

import (
	"time"

	"github.com/koykov/pbtk"
	"github.com/koykov/pbtk/frequency/cmsketch"
)

const (
	defaultDecayFactor     = .5
	defaultSoftDecayFactor = .75
)

type Config struct {
	cmsketch.Config
	// Count of added items to start decay.
	DecayLimit uint64
	// Time interval to start decay.
	DecayInterval time.Duration
	// External decay notifier to force decay start.
	ForceDecayNotifier ForceDecayNotifier
	// Default factor to decay counters.
	// Must be in range (0..1).
	DecayFactor float64
	// Soft factor to decay counters. Uses for too often decay operations.
	// Must be in range (0..1).
	SoftDecayFactor float64
}

func NewConfig(confidence, epsilon float64, hasher pbtk.Hasher) *Config {
	c := &Config{Config: cmsketch.Config{
		Concurrent: &cmsketch.ConcurrentConfig{}, // TinyLFU allows only concurrent CMS due to async decay
		Confidence: confidence,
		Epsilon:    epsilon,
		Hasher:     hasher,
	}}
	c.WithFlag(flagLFU, true)
	return c
}

func (c *Config) WithDecayLimit(limit uint64) *Config {
	c.DecayLimit = limit
	return c
}

func (c *Config) WithDecayInterval(interval time.Duration) *Config {
	c.DecayInterval = interval
	return c
}

func (c *Config) WithForceDecayNotifier(notifier ForceDecayNotifier) *Config {
	c.ForceDecayNotifier = notifier
	return c
}

func (c *Config) WithDecayFactor(df float64) *Config {
	c.DecayFactor = df
	return c
}

func (c *Config) WithSoftDecayFactor(sdf float64) *Config {
	c.SoftDecayFactor = sdf
	return c
}

func (c *Config) copy() *Config {
	cpy := *c
	return &cpy
}
