package frequency

import (
	"math"
	"testing"

	"github.com/koykov/approxity"
)

func TestMe[T []byte](t *testing.T, est Estimator[T]) {
	approxity.EachTestingDataset(func(_ int, ds *approxity.TestingDataset[[]byte]) {
		t.Run(ds.Name, func(t *testing.T) {
			est.Reset()
			for i := 0; i < len(ds.All); i++ {
				_ = est.Add(ds.All[i])
				if i != 0 && i%1000 == 0 {
					for j := 0; j < 1000; j++ {
						_ = est.Add(ds.All[i])
					}
				} else if i != 0 && i%100 == 0 {
					for j := 0; j < 100; j++ {
						_ = est.Add(ds.All[i])
					}
				} else if i != 0 && i%10 == 0 {
					for j := 0; j < 10; j++ {
						_ = est.Add(ds.All[i])
					}
				}
			}
			var diffv, diffc float64
			for i := 0; i < len(ds.All); i++ {
				var must uint64 = 1
				if i != 0 && i%1000 == 0 {
					must = 1001
				} else if i != 0 && i%100 == 0 {
					must = 101
				} else if i != 0 && i%10 == 0 {
					must = 11
				}
				e := est.Estimate(ds.All[i])
				if diff := math.Abs(float64(e) - float64(must)); diff > 0 {
					diffv += diff
					diffc++
				}
			}
			if diffc > 0 {
				t.Logf("avg diff: %f", diffv/diffc)
			}
		})
	})
}
