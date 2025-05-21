package amq

import (
	"context"
	"math"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/koykov/pbtk"
)

func TestMe[T []byte](t *testing.T, f Filter[T]) {
	pbtk.EachTestingDataset(func(_ int, ds *pbtk.TestingDataset[[]byte]) {
		t.Run(ds.Name, func(t *testing.T) {
			f.Reset()
			for i := 0; i < len(ds.Positives); i++ {
				_ = f.Set(ds.Positives[i])
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

func TestMeConcurrently[T []byte](t *testing.T, f Filter[T]) {
	pbtk.EachTestingDataset(func(_ int, ds *pbtk.TestingDataset[[]byte]) {
		if len(ds.All) == 0 {
			return
		}
		t.Run(ds.Name, func(t *testing.T) {
			f.Reset()
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			var wg sync.WaitGroup
			wg.Add(3)

			go func() {
				defer wg.Done()
				for i := 0; ; i++ {
					select {
					case <-ctx.Done():
						return
					default:
						_ = f.Set(ds.Positives[i%len(ds.Positives)])
					}
				}
			}()

			go func() {
				defer wg.Done()
				for i := 0; ; i++ {
					select {
					case <-ctx.Done():
						return
					default:
						_ = f.Unset(ds.All[i%len(ds.All)])
					}
				}
			}()

			go func() {
				defer wg.Done()
				for i := 0; ; i++ {
					select {
					case <-ctx.Done():
						return
					default:
						f.Contains(ds.All[(i % len(ds.All))])
					}
				}
			}()

			wg.Wait()
		})
	})
}

func BenchMe[T []byte](b *testing.B, f Filter[T]) {
	pbtk.EachTestingDataset(func(_ int, ds *pbtk.TestingDataset[[]byte]) {
		b.Run(ds.Name, func(b *testing.B) {
			f.Reset()
			for i := 0; i < len(ds.Positives); i++ {
				_ = f.Set(ds.Positives[i])
			}
			b.ReportAllocs()
			b.ResetTimer()
			for j := 0; j < b.N; j++ {
				f.Contains(ds.All[j%len(ds.All)])
			}
		})
	})
}

func BenchMeConcurrently[T []byte](b *testing.B, f Filter[T]) {
	pbtk.EachTestingDataset(func(_ int, ds *pbtk.TestingDataset[[]byte]) {
		b.Run(ds.Name, func(b *testing.B) {
			f.Reset()
			b.ReportAllocs()
			b.RunParallel(func(pb *testing.PB) {
				var i uint64 = math.MaxUint64
				for pb.Next() {
					ci := atomic.AddUint64(&i, 1)
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
