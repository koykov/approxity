package simhash

import (
	"github.com/koykov/byteseq"
	"github.com/koykov/pbtk"
	"github.com/koykov/pbtk/shingle"
)

type Config[T byteseq.Q] struct {
	Algo     pbtk.Hasher
	Shingler shingle.Shingler[T]
	K        uint
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
