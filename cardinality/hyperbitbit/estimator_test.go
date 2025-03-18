package hyperbitbit

import (
	"testing"

	"github.com/koykov/approxity/cardinality"
	"github.com/koykov/hash/xxhash"
)

const testN = 1e6

var testh = xxhash.Hasher64[[]byte]{}

func TestEstimator(t *testing.T) {
	t.Run("sync", func(t *testing.T) {
		est, err := NewEstimator[[]byte](NewConfig(testN, testh))
		if err != nil {
			t.Fatal(err)
		}
		cardinality.TestMe(t, est, -1) // disable delta checking due to HBB may be too inaccurate, especial on small datasets
	})
}

func BenchmarkEstimator(b *testing.B) {
	b.Run("sync", func(b *testing.B) {
		est, err := NewEstimator[[]byte](NewConfig(testN, testh))
		if err != nil {
			b.Fatal(err)
		}
		cardinality.BenchMe(b, est)
	})
}
