package cardinality

import (
	"context"
	"encoding/binary"
	"math"
	"sync"
	"sync/atomic"
	"testing"
	"time"
	"unsafe"

	"github.com/koykov/pbtk"
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
		if delta >= 0 && diff > delta {
			t.Errorf("estimation too inaccurate: ratio delta need %f, got %f", delta, diff)
		}
	})

	pbtk.EachTestingDataset(func(_ int, ds *pbtk.TestingDataset[[]byte]) {
		t.Run(ds.Name, func(t *testing.T) {
			est.Reset()
			for i := 0; i < len(ds.All); i++ {
				_ = est.Add(ds.All[i])
				if i%5 == 0 {
					// each 5th element adds 5 times
					for j := 0; j < 4; j++ {
						_ = est.Add(ds.All[i])
					}
				}
			}
			e := est.Estimate()
			ratio := float64(e) / float64(len(ds.All))
			diff := math.Abs(1 - ratio)
			if delta >= 0 && diff > delta {
				t.Errorf("estimation too inaccurate: ratio delta need %f, got %f", delta, diff)
			}
		})
	})
}

func TestMeConcurrently[T []byte](t *testing.T, est Estimator[T], delta float64) {
	t.Run("distinct counting", func(t *testing.T) {
		est.Reset()
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		var wg sync.WaitGroup

		for i := 0; i < 5; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				var buf [8]byte
				for j := 0; ; j++ {
					select {
					case <-ctx.Done():
						return
					default:
						binary.LittleEndian.PutUint64(buf[:], uint64(j))
						_ = est.Add(buf[:])
					}
				}
			}()
		}

		for i := 0; i < 5; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				tick := time.NewTicker(time.Millisecond * 5)
				defer tick.Stop()
				for {
					select {
					case <-ctx.Done():
						return
					case <-tick.C:
						est.Estimate()
					}
				}
			}()
		}

		wg.Wait()
	})
	pbtk.EachTestingDataset(func(_ int, ds *pbtk.TestingDataset[[]byte]) {
		t.Run(ds.Name, func(t *testing.T) {
			est.Reset()
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			var wg sync.WaitGroup
			wg.Add(2)

			go func() {
				defer wg.Done()
				for i := 0; ; i++ {
					select {
					case <-ctx.Done():
						return
					default:
						_ = est.Add(ds.All[i%len(ds.All)])
					}
				}
			}()

			go func() {
				defer wg.Done()
				tick := time.NewTicker(time.Millisecond * 5)
				defer tick.Stop()
				for i := 0; ; i++ {
					select {
					case <-ctx.Done():
						return
					case <-tick.C:
						est.Estimate()
					}
				}
			}()

			wg.Wait()

			e := est.Estimate()
			ratio := float64(e) / float64(len(ds.All))
			diff := math.Abs(1 - ratio)
			if diff > delta {
				t.Errorf("estimation too inaccurate: ratio delta need %f, got %f", delta, diff)
			}
		})
	})
}

func BenchMe(b *testing.B, est Estimator[[]byte]) {
	b.Run("add", func(b *testing.B) {
		b.ReportAllocs()
		est.Reset()
		var buf [8]byte
		for i := 0; i < b.N; i++ {
			binary.LittleEndian.PutUint64(buf[:], uint64(i))
			_ = est.Add(buf[:])
		}
	})
	b.Run("estimate", func(b *testing.B) {
		var buf [8]byte
		for i := uint64(0); i < 1e7; i++ {
			binary.LittleEndian.PutUint64(buf[:], i)
			_ = est.Add(buf[:])
		}
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = est.Estimate()
		}
	})

	pbtk.EachTestingDataset(func(_ int, ds *pbtk.TestingDataset[[]byte]) {
		b.Run(ds.Name, func(b *testing.B) {
			est.Reset()
			for i := 0; i < len(ds.All); i++ {
				_ = est.Add(ds.All[i])
			}
			b.ResetTimer()
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				est.Estimate()
			}
		})
	})
}

func BenchMeConcurrently[T []byte](b *testing.B, est Estimator[T]) {
	b.Run("add", func(b *testing.B) {
		b.ReportAllocs()
		est.Reset()
		var buf [8]byte
		for i := 0; i < b.N; i++ {
			binary.LittleEndian.PutUint64(buf[:], uint64(i))
			_ = est.Add(buf[:])
		}
	})
	b.Run("estimate", func(b *testing.B) {
		var buf [8]byte
		for i := uint64(0); i < 1e7; i++ {
			binary.LittleEndian.PutUint64(buf[:], i)
			_ = est.Add(buf[:])
		}
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = est.Estimate()
		}
	})

	pbtk.EachTestingDataset(func(_ int, ds *pbtk.TestingDataset[[]byte]) {
		b.Run(ds.Name, func(b *testing.B) {
			est.Reset()
			b.ReportAllocs()
			b.RunParallel(func(pb *testing.PB) {
				var i uint64 = math.MaxUint64
				for pb.Next() {
					ci := atomic.AddUint64(&i, 1)
					switch ci % 100 {
					case 99:
						est.Estimate()
					default:
						buf := *(*[8]byte)(unsafe.Pointer(&ci))
						_ = est.Add(buf[:])
					}
				}
			})
		})
	})
}
