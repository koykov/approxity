package similarity

import (
	"testing"

	"github.com/koykov/pbtk/simtest"
)

func TestMe(t *testing.T, est Estimator[[]byte], threshold float64) {
	simtest.EachTestingDataset(func(_ int, ds *simtest.Dataset) {
		t.Run(ds.Name, func(t *testing.T) {
			for j := 0; j < len(ds.Tuples); j++ {
				est.Reset()
				tp := &ds.Tuples[j]
				r, err := est.Estimate(tp.A, tp.B)
				if err != nil {
					t.Error(err)
				}
				if r > threshold {
					t.Errorf("threshold overflow: %f; '%s' vs '%s'", r, string(tp.A), string(tp.B))
				}
			}
		})
	})
}
