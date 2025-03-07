package quotient

import (
	"testing"

	"github.com/koykov/approxity"
	"github.com/koykov/hash/xxhash"
)

const (
	testSz  = 1e7
	testFPP = .01
)

var testh = xxhash.Hasher64[[]byte]{}

func TestFilter(t *testing.T) {
	f, err := NewFilter[[]byte](NewConfig(testSz, testFPP, testh))
	if err != nil {
		t.Fatal(err)
	}
	approxity.TestMe(t, f)
}

func BenchmarkFilter(b *testing.B) {
	f, err := NewFilter[[]byte](NewConfig(testSz, testFPP, testh))
	if err != nil {
		b.Fatal(err)
	}
	approxity.BenchMe(b, f)
}
