package xor

import "math"

func optimalSegmentLength(n, arity uint64) (r uint64) {
	switch {
	case arity == 3:
		r = 1 << uint64(math.Floor(math.Log(float64(n)))/math.Log(3.33)+2.11)
	case arity == 4:
		r = 1 << uint64(math.Floor(math.Log(float64(n)))/math.Log(2.91)+.5)
	default:
		r = math.MaxUint16 + 1
	}
	return
}

func optimalSizeFactor(n, arity uint64) (r float64) {
	switch {
	case arity == 3:
		r = math.Max(1.125, .875+.25*math.Log(1e6)/math.Log(float64(n)))
	case arity == 4:
		r = math.Max(1.075, .77+.305*math.Log(6e5)/math.Log(float64(n)))
	default:
		r = 2
	}
	return
}
