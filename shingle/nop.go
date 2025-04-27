package shingle

import "github.com/koykov/byteseq"

type nop[T byteseq.Q] struct{}

func NewNOP[T byteseq.Q]() Shingler[T] {
	return &nop[T]{}
}

func (sh *nop[T]) Shingle(s T) []T {
	var buf [1]T
	buf[0] = s
	return buf[:]
}

func (sh *nop[T]) AppendShingle(dst []T, s T) []T {
	dst = append(dst, s)
	return dst
}

func (sh *nop[T]) Each(s T, fn func(T)) {
	fn(s)
}

func (sh *nop[T]) Reset() {}

var _ = NewNOP[string]
