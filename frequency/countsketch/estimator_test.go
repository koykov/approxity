package countsketch

import (
	"testing"

	"github.com/koykov/approxity/frequency"
	"github.com/koykov/hash/xxhash"
)

const (
	testConfidence = 0.999
	testEpsilon    = 0.001
)

var testh = xxhash.Hasher64[[]byte]{}

func TestEstimator(t *testing.T) {
	t.Run("sync", func(t *testing.T) {
		est, err := NewEstimator[[]byte](NewConfig(testConfidence, testEpsilon, testh))
		if err != nil {
			t.Fatal(err)
		}
		frequency.TestMe(t, frequency.NewTestSignedAdapter(est))
	})
	t.Run("concurrent", func(t *testing.T) {
		est, err := NewEstimator[[]byte](NewConfig(testConfidence, testEpsilon, testh).
			WithConcurrency())
		if err != nil {
			t.Fatal(err)
		}
		frequency.TestMeConcurrently(t, frequency.NewTestSignedAdapter(est))
	})
}

func BenchmarkEstimator(b *testing.B) {
	b.Run("sync", func(b *testing.B) {
		est, err := NewEstimator[[]byte](NewConfig(testConfidence, testEpsilon, testh))
		if err != nil {
			b.Fatal(err)
		}
		frequency.BenchMe(b, frequency.NewTestSignedAdapter(est))
	})
	b.Run("concurrent", func(b *testing.B) {
		est, err := NewEstimator[[]byte](NewConfig(testConfidence, testEpsilon, testh).
			WithConcurrency().WithWriteAttemptsLimit(5))
		if err != nil {
			b.Fatal(err)
		}
		frequency.BenchMeConcurrently(b, frequency.NewTestSignedAdapter(est))
	})
}
