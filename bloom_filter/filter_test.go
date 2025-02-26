package bloom

import (
	"testing"

	"github.com/koykov/amq"
	"github.com/koykov/hash/xxhash"
)

const testsz = 1e7

var testh = xxhash.Hasher64[[]byte]{}

func TestFilter(t *testing.T) {
	t.Run("sync", func(t *testing.T) {
		f, err := NewFilter(NewConfig(testsz, testh).
			WithNumberHashFunctions(3))
		if err != nil {
			t.Fatal(err)
		}
		amq.TestMe(t, f)
	})
	t.Run("concurrent", func(t *testing.T) {
		f, err := NewFilter(NewConfig(testsz, testh).
			WithNumberHashFunctions(3).
			WithConcurrency().WithWriteAttemptsLimit(5))
		if err != nil {
			t.Fatal(err)
		}
		amq.TestMeConcurrently(t, f)
	})
}

func BenchmarkFilter(b *testing.B) {
	b.Run("sync", func(b *testing.B) {
		f, err := NewFilter(NewConfig(testsz, testh).
			WithNumberHashFunctions(3))
		if err != nil {
			b.Fatal(err)
		}
		amq.BenchMe(b, f)
	})
	b.Run("concurrent", func(b *testing.B) {
		f, err := NewFilter(NewConfig(testsz, testh).
			WithNumberHashFunctions(3).
			WithConcurrency().WithWriteAttemptsLimit(5))
		if err != nil {
			b.Fatal(err)
		}
		amq.BenchMeConcurrently(b, f)
	})
}
