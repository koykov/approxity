package hyperloglog

import "math"

var pow2 [math.MaxUint8]float64

func init() {
	for i := 0; i < math.MaxUint8; i++ {
		pow2[i] = math.Pow(2, float64(i))
	}
}
