package lsh

import "github.com/koykov/pbtk"

type Hasher[T pbtk.Hashable] interface {
	Add(value T) error
	Hash() []uint64
	AppendHash([]uint64) []uint64
	Reset()
}
