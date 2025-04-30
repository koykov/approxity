package cosine

import (
	"testing"

	"github.com/koykov/hash/xxhash"
	"github.com/koykov/pbtk/lsh/minhash"
	"github.com/koykov/pbtk/shingle"
	"github.com/koykov/pbtk/simtest"
)

var (
	testh       = xxhash.Hasher64[[]byte]{}
	testshc     = shingle.NewChar[[]byte](3, "") // 3-gram
	testshw     = shingle.NewWord[[]byte](2, "") // 2-word shingle
	testlshc, _ = minhash.NewHasher[[]byte](minhash.NewConfig[[]byte](testh, 50, testshc))
	testlshw, _ = minhash.NewHasher[[]byte](minhash.NewConfig[[]byte](testh, 50, testshw))
)

func TestEstimator(t *testing.T) {
	t.Run("char", func(t *testing.T) {
		simtest.EachTestingDataset(func(_ int, ds *simtest.Dataset) {
			t.Run(ds.Name, func(t *testing.T) {
				e, err := NewEstimator[[]byte](NewConfig[[]byte](testlshc))
				_ = err
				var r float64
				for j := 0; j < len(ds.Tuples); j++ {
					e.Reset()
					tp := &ds.Tuples[j]
					if r, err = e.Estimate(tp.A, tp.B); err != nil {
						t.Error(err)
					}
				}
				_ = r
			})
		})
	})
	t.Run("word", func(t *testing.T) {
		simtest.EachTestingDataset(func(_ int, ds *simtest.Dataset) {
			t.Run(ds.Name, func(t *testing.T) {
				e, err := NewEstimator[[]byte](NewConfig[[]byte](testlshw))
				_ = err
				var r float64
				for j := 0; j < len(ds.Tuples); j++ {
					e.Reset()
					tp := &ds.Tuples[j]
					if r, err = e.Estimate(tp.A, tp.B); err != nil {
						t.Error(err)
					}
				}
				_ = r
			})
		})
	})
}
