package countsketch

import "math"

func optimalWD(confidence, epsilon float64) (w, d uint64) {
	w, d = uint64(math.Ceil(1/math.Pow(epsilon, 2))), uint64(math.Ceil(math.Log(1/(1-confidence))))
	return
}
