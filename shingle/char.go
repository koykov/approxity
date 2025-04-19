package shingle

import (
	"github.com/koykov/byteseq"
)

type char[T byteseq.Q] struct {
	k uint
	c bool
}

func NewChar[T byteseq.Q](k uint, clean bool) Shingler[T] {
	return &char[T]{k, clean}
}

func (sh *char[T]) Shingle(s T) []T {
	buf := make([]T, 0, len(s)/int(sh.k))
	return sh.AppendShingle(buf, s)
}

func (sh *char[T]) AppendShingle(dst []T, s T) []T {
	// todo implement me
	return nil
}

func (sh *char[T]) Each(s T, fn func(T)) {
	// todo implement me
}

func (sh *char[T]) Reset() {
	// todo implement me
}
