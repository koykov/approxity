package hyperloglog

import "math"

var (
	// precalculated 1/2^n
	pow2d1 [math.MaxUint8]float64
	// precalculated non-zero term
	nzt [math.MaxUint8]float64
)

func init() {
	for i := 0; i < math.MaxUint8; i++ {
		pow2d1[i] = 1 / math.Pow(2, float64(i))
		nzt[i] = 1
	}
	nzt[0] = 0
}
