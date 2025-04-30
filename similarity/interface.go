package similarity

import "github.com/koykov/byteseq"

type Estimator[T byteseq.Q] interface {
	Estimate(a, b T) (float64, error)
	Reset()
}
