package dlcsketch

import (
	"math"
)

func optimalM(n, d uint64, cp float64) uint64 {
	delta := n / 10
	if delta == 0 {
		return 0
	}
	n_, d_ := float64(n), float64(d)
	m := n_ / (d_ * math.Log(1/cp))
	for i := 0; i < 100; i++ {
		a := -n_ / (m * d_)
		b := math.Exp(a)
		cp1 := math.Pow(1-b, d_)
		if cp1 <= cp {
			return uint64(m)
		}
		m += float64(delta)
	}
	return 0
}
