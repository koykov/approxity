package shingle

import "github.com/koykov/byteseq"

type nop[T byteseq.Q] struct{}

func NewNOP[T byteseq.Q]() Shingler[T] {
	return &nop[T]{}
}

func (c *nop[T]) Shingle(s T) []T {
	var buf [1]T
	buf[0] = s
	return buf[:]
}

func (c *nop[T]) AppendShingle(dst []T, s T) []T {
	dst = append(dst, s)
	return dst
}

func (c *nop[T]) Each(s T, fn func(T)) {
	fn(s)
}
