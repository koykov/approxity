package xor

import (
	"testing"

	"github.com/koykov/approxity"
)

func BenchmarkPool(b *testing.B) {
	b.Run("sync", func(b *testing.B) {
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
	b.Run("concurrent", func(b *testing.B) {
		approxity.EachTestingDataset(func(i int, ds *approxity.TestingDataset[[]byte]) {
			b.Run(ds.Name, func(b *testing.B) {
				b.ReportAllocs()
				b.RunParallel(func(pb *testing.PB) {
					for pb.Next() {
						f, _ := AcquireWithKeys[[]byte](&Config{Hasher: testh}, ds.Positives)
						Release(f)
					}
				})
			})
		})
	})
}
