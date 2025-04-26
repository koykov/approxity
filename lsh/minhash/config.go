package minhash

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
	N uint
	// Shingler to vector input data.
	// Mandatory param.
	Shingler shingle.Shingler[T]
	// One shingle size.
	// Mandatory param.
	K uint
}

func NewConfig[T byteseq.Q](algo pbtk.Hasher, n uint, shingler shingle.Shingler[T], k uint) *Config[T] {
	return &Config[T]{
		Algo:     algo,
		N:        n,
		Shingler: shingler,
		K:        k,
	}
}

func (c *Config[T]) copy() *Config[T] {
	cpy := *c
	return &cpy
}
