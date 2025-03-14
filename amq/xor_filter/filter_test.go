package xor

import (
	"math"
	"sync/atomic"
	"testing"

	"github.com/koykov/approxity"
	"github.com/koykov/hash/xxhash"
)

var testh = xxhash.Hasher64[[]byte]{}

func TestFilter(t *testing.T) {
	approxity.EachTestingDataset(func(i int, ds *approxity.TestingDataset[[]byte]) {
		t.Run(ds.Name, func(t *testing.T) {
			f, err := NewFilterWithKeys[[]byte](&Config{Hasher: testh}, ds.Positives)
			if err != nil {
				t.Fatal(err)
			}
			var falseNegative, falsePositive int
			for i := 0; i < len(ds.Negatives); i++ {
				if f.Contains(ds.Negatives[i]) {
					falsePositive++
				}
			}
			if falsePositive > 0 {
				// Just warn, it's OK to have small amount of false positives.
				t.Logf("%d of %d negatives (%d total) gives false positive value", falsePositive, len(ds.Negatives), len(ds.All))
			}
			for i := 0; i < len(ds.Positives); i++ {
				if !f.Contains(ds.Positives[i]) {
					falseNegative++
				}
			}
			if falseNegative > 0 {
				t.Errorf("%d of %d positives (%d total) gives false negative value", falseNegative, len(ds.Positives), len(ds.All))
			}
		})
	})
}

func BenchmarkFilter(b *testing.B) {
	b.Run("sync", func(b *testing.B) {
		approxity.EachTestingDataset(func(i int, ds *approxity.TestingDataset[[]byte]) {
			b.Run(ds.Name, func(b *testing.B) {
				f, _ := NewFilterWithKeys[[]byte](&Config{Hasher: testh}, ds.Positives)
				b.ReportAllocs()
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					f.Contains(ds.All[i%len(ds.All)])
				}
			})
		})
	})
	b.Run("concurrent", func(b *testing.B) {
		approxity.EachTestingDataset(func(i int, ds *approxity.TestingDataset[[]byte]) {
			b.Run(ds.Name, func(b *testing.B) {
				f, _ := NewFilterWithKeys[[]byte](&Config{Hasher: testh}, ds.Positives)
				b.ReportAllocs()
				b.RunParallel(func(pb *testing.PB) {
					var i uint64 = math.MaxUint64
					for pb.Next() {
						ci := atomic.AddUint64(&i, 1)
						f.Contains(ds.All[ci%uint64(len(ds.All))])
					}
				})
			})
		})
	})
	b.Run("pool", func(b *testing.B) {
		approxity.EachTestingDataset(func(i int, ds *approxity.TestingDataset[[]byte]) {
			b.Run(ds.Name, func(b *testing.B) {
				b.ReportAllocs()
				b.ResetTimer()
				for j := 0; j < b.N; j++ {
					f, _ := AcquireWithKeys[[]byte](&Config{Hasher: testh}, ds.Positives)
					Release(f)
				}
			})
		})
	})
}
