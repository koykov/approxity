package hyperloglog

import (
	"testing"

	"github.com/koykov/approxity/cardinality"
	"github.com/koykov/hash/xxhash"
)

const testP = 18

var testh = xxhash.Hasher64[[]byte]{}

func TestEstimator(t *testing.T) {
	t.Run("sync", func(t *testing.T) {
		est, err := NewEstimator[[]byte](NewConfig(testP, testh))
		if err != nil {
			t.Fatal(err)
		}
		cardinality.TestMe(t, est, 0.005)
	})
	t.Run("concurrent", func(t *testing.T) {
		est, err := NewEstimator[[]byte](NewConfig(testP, testh).
			WithConcurrency().WithWriteAttemptsLimit(5))
		if err != nil {
			t.Fatal(err)
		}
		cardinality.TestMeConcurrently(t, est, 0.005)
	})
}

func BenchmarkEstimator(b *testing.B) {
	b.Run("sync", func(b *testing.B) {
		est, err := NewEstimator[[]byte](NewConfig(testP, testh))
		if err != nil {
			b.Fatal(err)
		}
		cardinality.BenchMe(b, est)
	})
	b.Run("concurrent", func(b *testing.B) {
		est, err := NewEstimator[[]byte](NewConfig(testP, testh).
			WithConcurrency().WithWriteAttemptsLimit(5))
		if err != nil {
			b.Fatal(err)
		}
		cardinality.BenchMeConcurrently(b, est)
	})
}
