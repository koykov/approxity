package simhash

import (
	"github.com/koykov/byteseq"
	"github.com/koykov/pbtk"
	"github.com/koykov/pbtk/shingle"
)

type Config[T byteseq.Q] struct {
	// Hash algorithm to use.
	// Mandatory param.
	Algo pbtk.Hasher
	// Shingler to vector input data.
	// Mandatory param.
	Shingler shingle.Shingler[T]
	// Chunk size of vector item.
	// Mandatory param.
	K uint
}

func NewConfig[T byteseq.Q](algo pbtk.Hasher, shingler shingle.Shingler[T], k uint) *Config[T] {
	return &Config[T]{
		Algo:     algo,
		Shingler: shingler,
		K:        k,
	}
}

func (c *Config[T]) copy() *Config[T] {
	cpy := *c
	return &cpy
}
