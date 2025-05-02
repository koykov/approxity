package jaccard

import (
	"github.com/koykov/byteseq"
	"github.com/koykov/pbtk/lsh"
)

type Config[T byteseq.Q] struct {
	LSH lsh.Hasher[T]
}

func NewConfig[T byteseq.Q](lsh lsh.Hasher[T]) *Config[T] {
	return &Config[T]{LSH: lsh}
}

func (c *Config[T]) copy() *Config[T] {
	return &Config[T]{LSH: c.LSH}
}
