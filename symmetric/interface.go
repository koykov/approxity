package symmetric

import "github.com/koykov/byteseq"

type Differ[T byteseq.Q] interface {
	Diff(a, b T) (float64, error)
	Reset()
}
