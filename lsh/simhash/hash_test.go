package simhash

import (
	"testing"

	"github.com/koykov/hash/xxhash"
	"github.com/koykov/pbtk/lsh"
	"github.com/koykov/pbtk/shingle"
)

var (
	testh   = xxhash.Hasher64[[]byte]{}
	testshc = shingle.NewChar[[]byte](3, "") // 3-gram
	testshw = shingle.NewWord[[]byte](2, "") // 2-word shingle
)

func TestHash(t *testing.T) {
	t.Run("char", func(t *testing.T) {
		h, err := NewHasher[[]byte](NewConfig(testh, testshc))
		_ = err
		lsh.TestMe(t, h, lsh.TestDistHamming, 1, 1.0)
	})
	t.Run("word", func(t *testing.T) {
		h, err := NewHasher[[]byte](NewConfig(testh, testshw))
		_ = err
		lsh.TestMe(t, h, lsh.TestDistHamming, 1, 1.0)
	})
}

func BenchmarkHash(b *testing.B) {
	b.Run("char", func(b *testing.B) {
		h, err := NewHasher[[]byte](NewConfig(testh, testshc))
		_ = err
		lsh.BenchmarkMe(b, h)
	})
	b.Run("word", func(b *testing.B) {
		h, err := NewHasher[[]byte](NewConfig(testh, testshw))
		_ = err
		lsh.BenchmarkMe(b, h)
	})
}
