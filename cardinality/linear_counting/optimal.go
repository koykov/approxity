package linear

import "math"

func optimalM(n uint64, cp float64) uint64 {
	return uint64(math.Ceil(math.Abs(float64(n) / math.Log(1-cp))))
}
