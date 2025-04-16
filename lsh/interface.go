package lsh

import "github.com/koykov/pbtk"

type Hasher[T pbtk.Hashable] interface {
	Add(value T) error
	HAdd(hvalue uint64) error
	Hash() []uint64
	Reset()
}
