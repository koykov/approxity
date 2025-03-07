package hyperloglog

import (
	"testing"

	"github.com/koykov/approxity/cardinality"
	"github.com/koykov/hash/xxhash"
)

func TestEstimator(t *testing.T) {
	const p = 18
	est, err := NewEstimator[[]byte](&Config{Precision: p, Hasher: xxhash.Hasher64[[]byte]{}})
	if err != nil {
		t.Fatal(err)
	}
	cardinality.TestMe(t, est, 0.005)
}
