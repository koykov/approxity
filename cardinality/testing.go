package cardinality

import (
	"encoding/binary"
	"math"
	"testing"

	"github.com/koykov/approxity"
)

func TestMe[T []byte](t *testing.T, est Estimator[T], delta float64) {
	t.Run("distinct counting", func(t *testing.T) {
		const uniq uint64 = 1e7
		est.Reset()
		var buf [8]byte
		for i := 0; i < 10; i++ {
			for j := uint64(1); j < uniq; j++ {
				binary.LittleEndian.PutUint64(buf[:], j)
				_ = est.Add(buf[:])
			}
		}
		e := est.Estimate()
		ratio := float64(e) / float64(uniq)
		diff := math.Abs(1 - ratio)
		if diff > delta {
			t.Errorf("estimation too inaccurate: ratio delta need %f, got %f", delta, diff)
		}
	})

	approxity.EachTestingDataset(func(_ int, ds *approxity.TestingDataset[[]byte]) {
		t.Run(ds.Name, func(t *testing.T) {
			est.Reset()
			for i := 0; i < len(ds.All); i++ {
				_ = est.Add(ds.All[i])
				if i%5 == 0 {
					_ = est.Add(ds.All[i])
				}
			}
			e := est.Estimate()
			ratio := float64(e) / float64(len(ds.All))
			diff := math.Abs(1 - ratio)
			if diff > delta {
				t.Errorf("estimation too inaccurate: ratio delta need %f, got %f", delta, diff)
			}
		})
	})
}
