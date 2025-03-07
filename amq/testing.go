package amq

import (
	"context"
	"math"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/koykov/approxity"
)

func TestMe[T []byte](t *testing.T, f Filter[T]) {
	approxity.EachTestingDataset(func(i int, ds *approxity.TestingDataset[[]byte]) {
		t.Run(ds.Name, func(t *testing.T) {
			f.Reset()
			for j := 0; j < len(ds.Positives); j++ {
				_ = f.Set(ds.Positives[j])
			}
			var falseNegative, falsePositive int
			for j := 0; j < len(ds.Negatives); j++ {
				if f.Contains(ds.Negatives[j]) {
					falsePositive++
				}
			}
			if falsePositive > 0 {
				// Just warn, it's OK to have small amount of false positives.
				t.Logf("%d of %d negatives (%d total) gives false positive value", falsePositive, len(ds.Negatives), len(ds.All))
			}
			for j := 0; j < len(ds.Positives); j++ {
				if !f.Contains(ds.Positives[j]) {
					falseNegative++
				}
			}
			if falseNegative > 0 {
				t.Errorf("%d of %d positives (%d total) gives false Negatives value", falseNegative, len(ds.Positives), len(ds.All))
			}
		})
	})
}

func TestMeConcurrently[T []byte](t *testing.T, f Filter[T]) {
	approxity.EachTestingDataset(func(i int, ds *approxity.TestingDataset[[]byte]) {
		t.Run(ds.Name, func(t *testing.T) {
			f.Reset()
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			var wg sync.WaitGroup
			wg.Add(3)

			go func() {
				defer wg.Done()
				for j := 0; ; j++ {
					select {
					case <-ctx.Done():
						return
					default:
						_ = f.Set(ds.Positives[j%len(ds.Positives)])
					}
				}
			}()

			go func() {
				defer wg.Done()
				for j := 0; ; j++ {
					select {
					case <-ctx.Done():
						return
					default:
						_ = f.Unset(ds.All[j%len(ds.All)])
					}
				}
			}()

			go func() {
				defer wg.Done()
				for j := 0; ; j++ {
					select {
					case <-ctx.Done():
						return
					default:
						f.Contains(ds.All[(j % len(ds.All))])
					}
				}
			}()

			wg.Wait()
		})
	})
}

func BenchMe[T []byte](b *testing.B, f Filter[T]) {
	approxity.EachTestingDataset(func(i int, ds *approxity.TestingDataset[[]byte]) {
		b.Run(ds.Name, func(b *testing.B) {
			f.Reset()
			for j := 0; j < len(ds.Positives); j++ {
				_ = f.Set(ds.Positives[j])
			}
			b.ReportAllocs()
			b.ResetTimer()
			for k := 0; k < b.N; k++ {
				f.Contains(ds.All[k%len(ds.All)])
			}
		})
	})
}

func BenchMeConcurrently[T []byte](b *testing.B, f Filter[T]) {
	approxity.EachTestingDataset(func(i int, ds *approxity.TestingDataset[[]byte]) {
		b.Run(ds.Name, func(b *testing.B) {
			f.Reset()
			b.ReportAllocs()
			b.RunParallel(func(pb *testing.PB) {
				var j uint64 = math.MaxUint64
				for pb.Next() {
					ci := atomic.AddUint64(&j, 1)
					switch ci % 100 {
					case 99:
						_ = f.Set(ds.Positives[ci%uint64(len(ds.Positives))])
					case 98:
						_ = f.Unset(ds.All[ci%uint64(len(ds.All))])
					default:
						f.Contains(ds.All[ci%uint64(len(ds.All))])
					}
				}
			})
		})
	})
}
