package oddsketch

import "math"

// Calculate optimal filter size by given number of items (n) and false positive probability (fpp).
func optimalM(n uint64, fpp float64) uint64 {
	return uint64(math.Ceil(float64(-int64(n)) * math.Log(fpp) / (math.Pow(math.Log(2), 2))))
}

// Calculate optimal number of hash functions by given filter size (m) and number of items (n).
func optimalK(n, m uint64) uint64 {
	k := uint64(math.Ceil(math.Log(2) * float64(m) / float64(n)))
	if k == 0 {
		k = 1
	}
	return k
}
