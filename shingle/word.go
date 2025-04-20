package shingle

import (
	"github.com/koykov/byteconv"
	"github.com/koykov/byteseq"
)

type word[T byteseq.Q] struct {
	base[T]
	k uint
}

func NewWord[T byteseq.Q](k uint, cleanSet string) Shingler[T] {
	return &word[T]{
		base: base[T]{cset: byteconv.S2B(cleanSet)},
		k:    k,
	}
}

func (sh *word[T]) Shingle(s T) []T {
	// todo implement me
	return nil
}

func (sh *word[T]) AppendShingle(dts []T, s T) []T {
	// todo implement me
	return nil
}

func (sh *word[T]) Each(s T, fn func(T)) {
	// todo implement me
}

func (sh *word[T]) Reset() {
	// todo implement me
}
