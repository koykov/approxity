package spacesaving

import (
	"testing"

	"github.com/koykov/hash/xxhash"
	"github.com/koykov/pbtk/heavy"
)

var (
	testh     = xxhash.Hasher64[[]byte]{}
	testAlpha = 0.01
)

func TestHitter(t *testing.T) {
	h, err := NewHitter[[]byte](NewConfig(5, testh).
		WithEWMA(testAlpha))
	if err != nil {
		t.Fatal(err)
	}
	heavy.TestMe(t, h)
}

func BenchmarkHitter(b *testing.B) {
	h, err := NewHitter[[]byte](NewConfig(5, testh).
		WithEWMA(testAlpha))
	if err != nil {
		b.Fatal(err)
	}
	heavy.BenchMe(b, h)
}
