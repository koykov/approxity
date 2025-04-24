package shingle

import "github.com/koykov/byteseq"

type Shingler[T byteseq.Q] interface {
	Shingle(s T, k uint) []T
	AppendShingle(dst []T, s T, k uint) []T
	Each(s T, k uint, fn func(T))
	Reset()
}
