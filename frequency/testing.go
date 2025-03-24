package frequency

import (
	"context"
	"math"
	"sync"
	"testing"
	"time"

	"github.com/koykov/approxity"
)

func TestMe[T []byte](t *testing.T, est Estimator[T]) {
	approxity.EachTestingDataset(func(_ int, ds *approxity.TestingDataset[[]byte]) {
		t.Run(ds.Name, func(t *testing.T) {
			est.Reset()
			for i := 0; i < len(ds.All); i++ {
				_ = est.Add(ds.All[i])
				if i != 0 && i%1000 == 0 {
					for j := 0; j < 1000; j++ {
						_ = est.Add(ds.All[i])
					}
				} else if i != 0 && i%100 == 0 {
					for j := 0; j < 100; j++ {
						_ = est.Add(ds.All[i])
					}
				} else if i != 0 && i%10 == 0 {
					for j := 0; j < 10; j++ {
						_ = est.Add(ds.All[i])
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
				e := est.Estimate(ds.All[i])
				if diff := math.Abs(float64(e) - float64(must)); diff > 0 {
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

func TestMeConcurrently[T []byte](t *testing.T, est Estimator[T]) {
	approxity.EachTestingDataset(func(_ int, ds *approxity.TestingDataset[[]byte]) {
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
						est.Estimate(ds.All[i%len(ds.All)])
					}
				}
			}()

			wg.Wait()
		})
	})
}
