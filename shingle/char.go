package shingle

import "github.com/koykov/byteseq"

type char[T byteseq.Q] struct {
	k uint
	c bool
}

func NewChar[T byteseq.Q](k uint, clean bool) Shingler[T] {
	return &char[T]{k, clean}
}

func (c *char[T]) Shingle(s T) []T {
	// todo implement me
	return nil
}

func (c *char[T]) AppendShingle(dst []T, s T) []T {
	// todo implement me
	return nil
}

func (c *char[T]) Each(s T, fn func(T)) {
	// todo implement me
}
