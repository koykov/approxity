package countminsketch

import (
	"testing"

	"github.com/koykov/approxity/frequency"
	"github.com/koykov/hash/xxhash"
)

const (
	testConfidence = 0.99999
	testEpsilon    = 0.00001
)

var testh = xxhash.Hasher64[[]byte]{}

func TestEstimator(t *testing.T) {
	t.Run("sync", func(t *testing.T) {
		est, err := NewEstimator[[]byte](NewConfig(testConfidence, testEpsilon, testh))
		if err != nil {
			t.Fatal(err)
		}
		frequency.TestMe(t, est)
	})
	t.Run("concurrent", func(t *testing.T) {
		est, err := NewEstimator[[]byte](NewConfig(testConfidence, testEpsilon, testh).
			WithConcurrency())
		if err != nil {
			t.Fatal(err)
		}
		frequency.TestMeConcurrently(t, est)
	})
}
