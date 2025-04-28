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
	K uint64
	// Shingler to vector input data.
	// Mandatory param.
	Shingler shingle.Shingler[T]
	// Values storage.
	// If this param omitted, the instance of DefaultVector will be used.
	Vector Vector
}

func NewConfig[T byteseq.Q](algo pbtk.Hasher, k uint64, shingler shingle.Shingler[T]) *Config[T] {
	return &Config[T]{
		Algo:     algo,
		K:        k,
		Shingler: shingler,
	}
}

func (c *Config[T]) copy() *Config[T] {
	cpy := *c
	return &cpy
}
