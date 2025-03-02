package quotient

import "math"

// Calculate optimal filter size by given number of items (n),  false positive probability (fpp) and load factor (lf).
func optimalMQR(n uint64, fpp, lf float64) (m, q, r uint64) {
	q = uint64(math.Ceil(math.Log2(float64(n) / lf)))
	r = uint64(-math.Log2(fpp))
	b := (1 << q) * (r + 3)
	m = b / 8
	if b-m*8 > 0 {
		m++
	}
	return
}
