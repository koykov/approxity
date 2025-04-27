package shingle

import "github.com/koykov/byteseq"

type Shingler[T byteseq.Q] interface {
	Shingle(s T) []T
	AppendShingle(dst []T, s T) []T
	Each(s T, fn func(T))
	Reset()
}
