package shingle

import "github.com/koykov/byteseq"

type Shingler[T byteseq.Q] interface {
	Shingle(s T) []T
	AppendShingle(dts []T, s T) []T
	Each(fn func(T))
}
