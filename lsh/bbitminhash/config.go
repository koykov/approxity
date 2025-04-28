package bbitminhash

import (
	"github.com/koykov/byteseq"
	"github.com/koykov/pbtk"
	"github.com/koykov/pbtk/lsh/minhash"
	"github.com/koykov/pbtk/shingle"
)

type Config[T byteseq.Q] struct {
	minhash.Config[T]
	// Number of lower bits to store.
	// Mandatory param.
	B uint64
}

func NewConfig[T byteseq.Q](algo pbtk.Hasher, k uint64, shingler shingle.Shingler[T], b uint64) *Config[T] {
	return &Config[T]{
		Config: minhash.Config[T]{
			Algo:     algo,
			K:        k,
			Shingler: shingler,
			Vector:   newVector(b),
		},
		B: b,
	}
}

func (c *Config[T]) copy() *Config[T] {
	cpy := *c
	return &cpy
}
