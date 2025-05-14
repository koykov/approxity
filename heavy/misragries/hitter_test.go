package misragries

import (
	"testing"

	"github.com/koykov/hash/xxhash"
	"github.com/koykov/pbtk/heavy"
)

var testh = xxhash.Hasher64[[]byte]{}

func TestHitter(t *testing.T) {
	h, err := NewHitter[[]byte](NewConfig(5, testh).WithBuckets(1))
	if err != nil {
		t.Fatal(err)
	}
	heavy.TestMe(t, h, 20)
}

func BenchmarkHitter(b *testing.B) {
	h, err := NewHitter[[]byte](NewConfig(5, testh))
	if err != nil {
		b.Fatal(err)
	}
	heavy.BenchMe(b, h)
}
