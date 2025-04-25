package simhash

import (
	"testing"

	"github.com/koykov/hash/xxhash"
	"github.com/koykov/pbtk/lsh"
	"github.com/koykov/pbtk/shingle"
)

var (
	testh  = xxhash.Hasher64[[]byte]{}
	testsh = shingle.NewChar[[]byte]("")
	testk  = uint(3)
)

func TestHash(t *testing.T) {
	h, err := NewHasher[[]byte](NewConfig(testh, testsh, testk))
	_ = err
	lsh.TestMe(t, h, lsh.TestDistHamming, 1.0)
}

func BenchmarkHash(b *testing.B) {
	h, err := NewHasher[[]byte](NewConfig(testh, testsh, testk))
	_ = err
	lsh.BenchmarkMe(b, h)
}
