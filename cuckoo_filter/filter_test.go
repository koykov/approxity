package cuckoo

import (
	"testing"

	"github.com/koykov/amq"
	"github.com/koykov/hash/metro"
)

const (
	testsz   = 1e7
	testseed = uint64(1234)
)

var testh = metro.Hasher64[[]byte]{Seed: testseed}

func TestFilter(t *testing.T) {
	t.Run("sync", func(t *testing.T) {
		f, err := NewFilter(NewConfig(testsz, testh))
		if err != nil {
			t.Fatal(err)
		}
		amq.TestMe(t, f)
	})
	t.Run("concurrent", func(t *testing.T) {
		f, err := NewFilter(NewConfig(testsz, testh).
			WithConcurrency().WithWriteAttemptsLimit(5))
		if err != nil {
			t.Fatal(err)
		}
		amq.TestMeConcurrently(t, f)
	})
}

func BenchmarkFilter(b *testing.B) {
	b.Run("sync", func(b *testing.B) {
		f, err := NewFilter(NewConfig(testsz, testh))
		if err != nil {
			b.Fatal(err)
		}
		amq.BenchMe(b, f)
	})
	b.Run("concurrent", func(b *testing.B) {
		f, err := NewFilter(NewConfig(testsz, testh).
			WithConcurrency().WithWriteAttemptsLimit(5))
		if err != nil {
			b.Fatal(err)
		}
		amq.BenchMeConcurrently(b, f)
	})
}
