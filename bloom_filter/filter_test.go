package bloom

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
	t.Run("sync", func(t *testing.T) {
		f, err := NewFilter(NewConfig(testSz, testFPP, testh))
		if err != nil {
			t.Fatal(err)
		}
		amq.TestMe(t, f)
	})
	t.Run("concurrent", func(t *testing.T) {
		f, err := NewFilter(NewConfig(testSz, testFPP, testh).
			WithConcurrency().WithWriteAttemptsLimit(5))
		if err != nil {
			t.Fatal(err)
		}
		amq.TestMeConcurrently(t, f)
	})
}

func BenchmarkFilter(b *testing.B) {
	b.Run("sync", func(b *testing.B) {
		f, err := NewFilter(NewConfig(testSz, testFPP, testh))
		if err != nil {
			b.Fatal(err)
		}
		amq.BenchMe(b, f)
	})
	b.Run("concurrent", func(b *testing.B) {
		f, err := NewFilter(NewConfig(testSz, testFPP, testh).
			WithConcurrency().WithWriteAttemptsLimit(5))
		if err != nil {
			b.Fatal(err)
		}
		amq.BenchMeConcurrently(b, f)
	})
}
