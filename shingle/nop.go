package shingle

import "github.com/koykov/byteseq"

type nop[T byteseq.Q] struct{}

func NewNOP[T byteseq.Q]() Shingler[T] {
	return &nop[T]{}
}

func (sh *nop[T]) Shingle(s T, _ uint) []T {
	var buf [1]T
	buf[0] = s
	return buf[:]
}

func (sh *nop[T]) AppendShingle(dst []T, s T, _ uint) []T {
	dst = append(dst, s)
	return dst
}

func (sh *nop[T]) Each(s T, _ uint, fn func(T)) {
	fn(s)
}

func (sh *nop[T]) Reset() {}
