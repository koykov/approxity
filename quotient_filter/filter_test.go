package quotient

import (
	"testing"

	"github.com/koykov/amq"
	"github.com/koykov/hash/xxhash"
)

const (
	testSz  = 1e6
	testFPP = .01
)

var testh = xxhash.Hasher64[[]byte]{}

func TestFilter(t *testing.T) {
	f, err := NewFilter(NewConfig(testSz, testFPP, testh))
	if err != nil {
		t.Fatal(err)
	}
	amq.TestMe(t, f)
}

func BenchmarkFilter(b *testing.B) {
	f, err := NewFilter(NewConfig(testSz, testFPP, testh))
	if err != nil {
		b.Fatal(err)
	}
	amq.BenchMe(b, f)
}
