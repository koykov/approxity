package spacesaving

import (
	"math/rand"
	"testing"

	"github.com/koykov/hash/xxhash"
	"github.com/koykov/pbtk"
	"github.com/koykov/pbtk/heavy"
)

var (
	testh     = xxhash.Hasher64[[]byte]{}
	testAlpha = 0.01
)

func TestHitter(t *testing.T) {
	pbtk.EachTestingDataset(func(_ int, ds *pbtk.TestingDataset[[]byte]) {
		t.Run(ds.Name, func(t *testing.T) {
			if len(ds.All) <= 20 {
				return
			}
			h, err := NewHitter[[]byte](NewConfig(5, testh).
				WithEWMA(testAlpha))
			if err != nil {
				t.Fatal(err)
			}
			repeat := len(ds.All) / 20
			for i := 0; i < len(ds.All); i++ {
				if err = h.Add(ds.All[i]); err != nil {
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
			t.Logf("%#v", r)
		})
	})
}

func BenchmarkHitter(b *testing.B) {
	pbtk.EachTestingDataset(func(_ int, ds *pbtk.TestingDataset[[]byte]) {
		b.Run(ds.Name, func(b *testing.B) {
			h, err := NewHitter[[]byte](NewConfig(5, testh).
				WithEWMA(testAlpha))
			if err != nil {
				b.Fatal(err)
			}
			b.Run("add", func(b *testing.B) {
				b.ReportAllocs()
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					if err = h.Add(ds.All[i%len(ds.All)]); err != nil {
						b.Fatal(err)
					}
				}
			})
			b.Run("hits", func(b *testing.B) {
				h.Reset()
				for i := 0; i < len(ds.All); i++ {
					if err = h.Add(ds.All[i]); err != nil {
						b.Fatal(err)
					}
				}
				b.ReportAllocs()
				b.ResetTimer()
				var buf []heavy.Hit[[]byte]
				for i := 0; i < b.N; i++ {
					buf = h.AppendHits(buf[:0])
				}
				_ = buf
			})
		})
	})
}
