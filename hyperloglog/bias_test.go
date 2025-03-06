package hyperloglog

import "testing"

var biasStages = []struct {
	p uint64
	e float64
	r float64
}{
	{14, 189094.71188332525, 185139.77039464746},
}

func TestBias(t *testing.T) {
	fuzzeq := func(a, b, e float64) bool {
		return a-e <= b && a+e >= b
	}
	for i := 0; i < len(biasStages); i++ {
		stage := &biasStages[i]
		t.Run("", func(t *testing.T) {
			r := biasEstimation(stage.p, stage.e)
			if !fuzzeq(r, stage.r, 0.0001) {
				t.Errorf("biasEstimation(%d, %f) = %f, want %f", stage.p, stage.r, r, stage.r)
			}
		})
	}
}

func BenchmarkBias(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for j := 0; j < len(biasStages); j++ {
			stage := &biasStages[j]
			b.Run("", func(b *testing.B) {
				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					_ = biasEstimation(stage.p, stage.e)
				}
			})
		}
	}
}
