package bbitminhash

import (
	"github.com/koykov/byteseq"
	"github.com/koykov/pbtk"
	"github.com/koykov/pbtk/shingle"
)

type Config[T byteseq.Q] struct {
	// Hash algorithm to use.
	// Mandatory param.
	Algo pbtk.Hasher
	// Number of hash functions.
	// Mandatory param.
	K uint64
	// Shingler to vector input data.
	// Mandatory param.
	Shingler shingle.Shingler[T]
	// Number of lower bits to store.
	// Mandatory param.
	B uint64
}

func NewConfig[T byteseq.Q](algo pbtk.Hasher, k uint64, shingler shingle.Shingler[T], b uint64) *Config[T] {
	return &Config[T]{
		Algo:     algo,
		K:        k,
		Shingler: shingler,
		B:        b,
	}
}

func (c *Config[T]) copy() *Config[T] {
	cpy := *c
	return &cpy
}
