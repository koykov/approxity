package oddsketch

import (
	"github.com/koykov/byteseq"
	"github.com/koykov/pbtk/lsh"
)

const defaultFPP = .01

type Config[T byteseq.Q] struct {
	// Number of desired items to store in the filter
	// Mandatory param.
	ItemsNumber uint64
	// False positive probability value.
	// If this param omit, defaultFPP (0.01) will use instead.
	FPP float64
	// Hasher to calculate hash sum of the items.
	// Mandatory param.
	LSH lsh.Hasher[T]
	// Setting up this section enables concurrent read/write operations.
	Concurrent *ConcurrentConfig
}

// ConcurrentConfig configures concurrent section of config.
type ConcurrentConfig struct {
	// How many write attempts may perform.
	WriteAttemptsLimit uint64
}

func NewConfig[T byteseq.Q](items uint64, fpp float64, lsh lsh.Hasher[T]) *Config[T] {
	return &Config[T]{
		ItemsNumber: items,
		FPP:         fpp,
		LSH:         lsh,
	}
}

func (c *Config[T]) WithConcurrency() *Config[T] {
	c.Concurrent = &ConcurrentConfig{}
	return c
}

func (c *Config[T]) WithItemsNumber(items uint64) *Config[T] {
	c.ItemsNumber = items
	return c
}

func (c *Config[T]) WithFPP(fpp float64) *Config[T] {
	c.FPP = fpp
	return c
}

func (c *Config[T]) WithLSH(lsh lsh.Hasher[T]) *Config[T] {
	c.LSH = lsh
	return c
}

func (c *Config[T]) WithWriteAttemptsLimit(limit uint64) *Config[T] {
	if c.Concurrent == nil {
		c.Concurrent = &ConcurrentConfig{}
	}
	c.Concurrent.WriteAttemptsLimit = limit
	return c
}

func (c *Config[T]) copy() *Config[T] {
	cpy := *c
	return &cpy
}
