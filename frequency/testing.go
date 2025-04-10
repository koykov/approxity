package frequency

import (
	"context"
	"encoding/binary"
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/koykov/pbtk"
)

type TestAdapter[T []byte] struct {
	signed   SignedEstimator[T]
	unsigned Estimator[T]
	precise  PreciseEstimator[T]
}

func NewTestAdapter[T []byte](est Estimator[T]) *TestAdapter[T] {
	return &TestAdapter[T]{unsigned: est}
}

func NewTestSignedAdapter[T []byte](est SignedEstimator[T]) *TestAdapter[T] {
	return &TestAdapter[T]{signed: est}
}

func NewTestPreciseAdapter[T []byte](est PreciseEstimator[T]) *TestAdapter[T] {
	return &TestAdapter[T]{precise: est}
}

func (t *TestAdapter[T]) Add(key T) error {
	switch {
	case t.unsigned != nil:
		return t.unsigned.Add(key)
	case t.signed != nil:
		return t.signed.Add(key)
	case t.precise != nil:
		return t.precise.Add(key)
	}
	return fmt.Errorf("no estimator found")
}

func (t *TestAdapter[T]) HAdd(hkey uint64) error {
	switch {
	case t.unsigned != nil:
		return t.unsigned.HAdd(hkey)
	case t.signed != nil:
		return t.signed.HAdd(hkey)
	case t.precise != nil:
		return t.precise.HAdd(hkey)
	}
	return fmt.Errorf("no estimator found")
}

func (t *TestAdapter[T]) Reset() {
	switch {
	case t.unsigned != nil:
		t.unsigned.Reset()
	case t.signed != nil:
		t.signed.Reset()
	}
}

func (t *TestAdapter[T]) StubEstimate(key T) {
	switch {
	case t.unsigned != nil:
		t.unsigned.Estimate(key)
	case t.signed != nil:
		t.signed.Estimate(key)
	case t.precise != nil:
		t.precise.Estimate(key)
	}
}

func TestMe[T []byte](t *testing.T, a *TestAdapter[T]) {
	pbtk.EachTestingDataset(func(_ int, ds *pbtk.TestingDataset[[]byte]) {
		t.Run(ds.Name, func(t *testing.T) {
			a.Reset()
			for i := 0; i < len(ds.All); i++ {
				_ = a.Add(ds.All[i])
				if i != 0 && i%1000 == 0 {
					for j := 0; j < 1000; j++ {
						_ = a.Add(ds.All[i])
					}
				} else if i != 0 && i%100 == 0 {
					for j := 0; j < 100; j++ {
						_ = a.Add(ds.All[i])
					}
				} else if i != 0 && i%10 == 0 {
					for j := 0; j < 10; j++ {
						_ = a.Add(ds.All[i])
					}
				}
			}
			var diffv, diffc float64
			for i := 0; i < len(ds.All); i++ {
				var must uint64 = 1
				if i != 0 && i%1000 == 0 {
					must = 1001
				} else if i != 0 && i%100 == 0 {
					must = 101
				} else if i != 0 && i%10 == 0 {
					must = 11
				}
				var e float64
				switch {
				case a.unsigned != nil:
					e = float64(a.unsigned.Estimate(ds.All[i]))
				case a.signed != nil:
					e = float64(a.signed.Estimate(ds.All[i]))
				case a.precise != nil:
					e = a.precise.Estimate(ds.All[i])
				}
				if diff := math.Abs(e - float64(must)); diff > 0 {
					diffv += diff
					diffc++
				}
			}
			if diffc > 0 {
				t.Logf("avg diff: %f", diffv/diffc)
			}
		})
	})
}

func TestMeConcurrently[T []byte](t *testing.T, a *TestAdapter[T]) {
	pbtk.EachTestingDataset(func(_ int, ds *pbtk.TestingDataset[[]byte]) {
		t.Run(ds.Name, func(t *testing.T) {
			a.Reset()
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
						_ = a.Add(ds.All[i%len(ds.All)])
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
						a.StubEstimate(ds.All[i%len(ds.All)])
					}
				}
			}()

			wg.Wait()
		})
	})
}

func BenchMe[T []byte](b *testing.B, a *TestAdapter[T]) {
	b.Run("add", func(b *testing.B) {
		b.ReportAllocs()
		a.Reset()
		var buf [8]byte
		for i := 0; i < b.N; i++ {
			binary.LittleEndian.PutUint64(buf[:], uint64(i))
			_ = a.Add(buf[:])
		}
	})
	b.Run("estimate", func(b *testing.B) {
		a.Reset()
		var buf [8]byte
		for i := uint64(0); i < 1e7; i++ {
			binary.LittleEndian.PutUint64(buf[:], i)
			_ = a.Add(buf[:])
		}
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			binary.LittleEndian.PutUint64(buf[:], uint64(i))
			a.StubEstimate(buf[:])
		}
	})
	pbtk.EachTestingDataset(func(_ int, ds *pbtk.TestingDataset[[]byte]) {
		b.Run(ds.Name, func(b *testing.B) {
			a.Reset()
			for i := 0; i < len(ds.All); i++ {
				_ = a.Add(ds.All[i])
			}
			b.ResetTimer()
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				a.StubEstimate(ds.All[i%len(ds.All)])
			}
		})
	})
}

func BenchMeConcurrently[T []byte](b *testing.B, a *TestAdapter[T]) {
	pbtk.EachTestingDataset(func(_ int, ds *pbtk.TestingDataset[[]byte]) {
		b.Run(ds.Name, func(b *testing.B) {
			var pool = sync.Pool{New: func() any {
				var buf [8]byte
				return &buf
			}}
			a.Reset()
			b.ReportAllocs()
			b.RunParallel(func(pb *testing.PB) {
				var i uint64 = math.MaxUint64
				for pb.Next() {
					ci := atomic.AddUint64(&i, 1)
					buf := pool.Get().(*[8]byte)
					binary.LittleEndian.PutUint64((*buf)[:], ci)
					switch ci % 100 {
					case 99:
						a.StubEstimate(buf[:])
					default:
						_ = a.Add(buf[:])
					}
					pool.Put(buf)
				}
			})
		})
	})
}
