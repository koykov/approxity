package heavy

import (
	"math/rand"
	"testing"

	"github.com/koykov/pbtk"
)

func TestMe[T []byte](t *testing.T, h Hitter[T], repeatRange int) {
	pbtk.EachTestingDataset(func(_ int, ds *pbtk.TestingDataset[[]byte]) {
		t.Run(ds.Name, func(t *testing.T) {
			h.Reset()
			if len(ds.All) <= repeatRange {
				return
			}
			repeat := len(ds.All) / repeatRange
			for i := 0; i < len(ds.All); i++ {
				if err := h.Add(ds.All[i]); err != nil {
					t.Fatal(err)
				}
				if i%repeat == 0 {
					j := rand.Intn(6)
					_ = h.Add(ds.All[j])
				}
			}
			hits := h.Hits()
			r := make(map[string]float64, len(hits))
			for i := 0; i < len(hits); i++ {
				r[string(hits[i].Key)] = hits[i].Rate
			}
			t.Logf("%v", r)
		})
	})
}

func BenchMe(b *testing.B, h Hitter[[]byte]) {
	pbtk.EachTestingDataset(func(_ int, ds *pbtk.TestingDataset[[]byte]) {
		b.Run(ds.Name, func(b *testing.B) {
			b.Run("add", func(b *testing.B) {
				h.Reset()
				b.ReportAllocs()
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					if err := h.Add(ds.All[i%len(ds.All)]); err != nil {
						b.Fatal(err)
					}
				}
			})
			b.Run("hits", func(b *testing.B) {
				h.Reset()
				for i := 0; i < len(ds.All); i++ {
					if err := h.Add(ds.All[i]); err != nil {
						b.Fatal(err)
					}
				}
				b.ReportAllocs()
				b.ResetTimer()
				var buf []Hit[[]byte]
				for i := 0; i < b.N; i++ {
					buf = h.AppendHits(buf[:0])
				}
				_ = buf
			})
		})
	})
}
