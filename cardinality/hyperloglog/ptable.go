package hyperloglog

import "math"

var (
	// precalculated 1/2^n
	pow2d1 [math.MaxUint8]float64
	// precalculated non-zero term
	nzt [math.MaxUint8]float64
	// computed threshold for each possible precision
	threshold = [15]float64{10, 20, 40, 80, 220, 400, 900, 1800, 3100, 6500, 11500, 20000, 50000, 120000, 350000}
)

func init() {
	for i := 0; i < math.MaxUint8; i++ {
		pow2d1[i] = 1 / math.Pow(2, float64(i))
		nzt[i] = 1
	}
	nzt[0] = 0
}
