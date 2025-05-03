package symmetric

import (
	"testing"

	"github.com/koykov/pbtk/simtest"
)

func TestMe(t *testing.T, diff Differ[[]byte], threshold float64) {
	simtest.EachTestingDataset(func(_ int, ds *simtest.Dataset) {
		t.Run(ds.Name, func(t *testing.T) {
			for j := 0; j < len(ds.Tuples); j++ {
				diff.Reset()
				tp := &ds.Tuples[j]
				r, err := diff.Diff(tp.A, tp.B)
				if err != nil {
					t.Error(err)
				}
				if r < threshold {
					t.Errorf("diff = %f, expected >= %f", r, threshold)
				}
			}
		})
	})
}

func BenchMe(b *testing.B, diff Differ[[]byte]) {
	simtest.EachTestingDataset(func(_ int, ds *simtest.Dataset) {
		b.Run(ds.Name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				for j := 0; j < len(ds.Tuples); j++ {
					diff.Reset()
					tp := &ds.Tuples[j]
					_, _ = diff.Diff(tp.A, tp.B)
				}
			}
		})
	})
}
