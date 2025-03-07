package hyperloglog

import (
	"testing"

	"github.com/koykov/approxity/cardinality"
	"github.com/koykov/hash/xxhash"
)

const testP = 18

var testh = xxhash.Hasher64[[]byte]{}

func TestEstimator(t *testing.T) {
	est, err := NewEstimator[[]byte](&Config{Precision: testP, Hasher: testh})
	if err != nil {
		t.Fatal(err)
	}
	cardinality.TestMe(t, est, 0.005)
}

func BenchmarkEstimator(b *testing.B) {
	est, err := NewEstimator[[]byte](&Config{Precision: testP, Hasher: testh})
	if err != nil {
		b.Fatal(err)
	}
	cardinality.BenchMe(b, est)
}
