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
}
