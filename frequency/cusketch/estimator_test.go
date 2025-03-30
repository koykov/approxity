package cusketch

import (
	"testing"

	"github.com/koykov/hash/xxhash"
	"github.com/koykov/pbtk/frequency"
	"github.com/koykov/pbtk/frequency/cmsketch"
)

const (
	testConfidence = 0.99999
	testEpsilon    = 0.00001
)

var testh = xxhash.Hasher64[[]byte]{}

func TestEstimator(t *testing.T) {
	t.Run("sync", func(t *testing.T) {
		est, err := NewEstimator[[]byte](cmsketch.NewConfig(testConfidence, testEpsilon, testh))
		if err != nil {
			t.Fatal(err)
		}
		frequency.TestMe(t, frequency.NewTestAdapter(est))
	})
	t.Run("concurrent", func(t *testing.T) {
		est, err := NewEstimator[[]byte](cmsketch.NewConfig(testConfidence, testEpsilon, testh).
			WithConcurrency())
		if err != nil {
			t.Fatal(err)
		}
		frequency.TestMeConcurrently(t, frequency.NewTestAdapter(est))
	})
}

func BenchmarkEstimator(b *testing.B) {
	b.Run("sync", func(b *testing.B) {
		est, err := NewEstimator[[]byte](cmsketch.NewConfig(testConfidence, testEpsilon, testh))
		if err != nil {
			b.Fatal(err)
		}
		frequency.BenchMe(b, frequency.NewTestAdapter(est))
	})
	b.Run("concurrent", func(b *testing.B) {
		est, err := NewEstimator[[]byte](cmsketch.NewConfig(testConfidence, testEpsilon, testh).
			WithConcurrency().WithWriteAttemptsLimit(5))
		if err != nil {
			b.Fatal(err)
		}
		frequency.BenchMeConcurrently(b, frequency.NewTestAdapter(est))
	})
}
