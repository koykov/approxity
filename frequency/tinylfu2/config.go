package tinylfu

import (
	"github.com/koykov/pbtk"
	"github.com/koykov/pbtk/frequency"
)

const (
	defaultTau          = 30 // 30 seconds
	defaultMinDeltaTime = 1
)

type Config struct {
	// Confidence represent a possibility that potential error will be in range of acceptable error rate (see Epsilon).
	// E.g.: Confidence 0.99 guarantees 99% of estimations will be in range represent by Epsilon.
	// Must be in range (0..1).
	// Mandatory param.
	Confidence float64
	// Epsilon represent precision of estimation: less epsilon value makes estimation more accurate, but grows the table.
	// E.g.: Epsilon 0.01 guarantees that estimation error is less or equal 1% of all elements.
	// Must be in range (0..1).
	// Mandatory param.
	Epsilon float64
	// Hasher to calculate hash sum of the items.
	// Mandatory param.
	Hasher pbtk.Hasher
	// EWMA settings.
	EWMA EWMA
	// Clock to measure time deltas. Testing stuff.
	Clock Clock
	// Setting up this section enables concurrent read/write operations.
	Concurrent *ConcurrentConfig
	// Metrics writer handler.
	MetricsWriter frequency.MetricsWriter
}

// EWMA (exponentially weighted moving average) params.
type EWMA struct {
	// Smoothing constant (time in seconds) to decay Count-Min Sketch counters.
	Tau uint64
	// Minimal time delta to apply native EWMA.
	// For less time deltas uses hybrid approach - sum of old value with EWMA (e^(-MinDeltaTime/Tau)).
	// Hybrid approach allows to handle quick updates and keep precision/stability balance of EWMA.
	MinDeltaTime uint64
}

// ConcurrentConfig configures concurrent section of config.
type ConcurrentConfig struct {
	// How many write attempts may perform.
	WriteAttemptsLimit uint64
}

func NewConfig(confidence, epsilon float64, hasher pbtk.Hasher) *Config {
	return &Config{
		Confidence: confidence,
		Epsilon:    epsilon,
		Hasher:     hasher,
	}
}

func (c *Config) WithConcurrency() *Config {
	c.Concurrent = &ConcurrentConfig{}
	return c
}

func (c *Config) WithWriteAttemptsLimit(limit uint64) *Config {
	if c.Concurrent == nil {
		c.Concurrent = &ConcurrentConfig{}
	}
	c.Concurrent.WriteAttemptsLimit = limit
	return c
}

func (c *Config) WithEWMA(tau, minDeltaTime uint64) *Config {
	c.EWMA.Tau = tau
	c.EWMA.MinDeltaTime = minDeltaTime
	return c
}

func (c *Config) WithEWMATau(tau uint64) *Config {
	c.EWMA.Tau = tau
	return c
}

func (c *Config) WithEWMAminDeltaTime(minDeltaTime uint64) *Config {
	c.EWMA.MinDeltaTime = minDeltaTime
	return c
}

func (c *Config) WithClock(clock Clock) *Config {
	c.Clock = clock
	return c
}

func WithMetricsWriter(conf *Config, mw frequency.MetricsWriter) *Config {
	conf.MetricsWriter = mw
	return conf
}

func (c *Config) copy() *Config {
	cpy := *c
	return &cpy
}
