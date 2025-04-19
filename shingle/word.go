package shingle

import "github.com/koykov/byteseq"

type word[T byteseq.Q] struct {
	k uint
	c bool
}

func NewWord[T byteseq.Q](k uint, clean bool) Shingler[T] {
	return &word[T]{k, clean}
}

func (w *word[T]) Shingle(s T) []T {
	// todo implement me
	return nil
}

func (w *word[T]) AppendShingle(dts []T, s T) []T {
	// todo implement me
	return nil
}

func (w *word[T]) Each(fn func(T)) {
	// todo implement me
}
