package tinylfu

import "math"

func optimalWD(confidence, epsilon float64) (w, d uint64) {
	w, d = uint64(math.Ceil(math.E/epsilon)), uint64(math.Ceil(math.Log(1/(1-confidence))))
	return
}
