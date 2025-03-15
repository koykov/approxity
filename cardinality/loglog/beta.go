package loglog

import "math"

// beta-function correction coefficients
var beta = [8]uint64{
	0xbfd7b488a9987b73, 0x3fb20a70ff131279, 0x3fc6439022a26c51, 0x3fc4ea3d0aa27058,
	0xbfb7a60c6ea34bca, 0x3fa32381ba54d011, 0xbf760db32d24b781, 0x3f3bccba2d62e4a7,
}

func betaEstimation(z float64) float64 {
	zl := math.Log(z + 1)
	_ = beta[7]
	return math.Float64frombits(beta[0])*z + math.Float64frombits(beta[1])*zl +
		math.Float64frombits(beta[2])*math.Pow(zl, 2) + math.Float64frombits(beta[3])*math.Pow(zl, 3) +
		math.Float64frombits(beta[4])*math.Pow(zl, 4) + math.Float64frombits(beta[5])*math.Pow(zl, 5) +
		math.Float64frombits(beta[6])*math.Pow(zl, 6) + math.Float64frombits(beta[7])*math.Pow(zl, 7)
}
