package simhash

import (
	"testing"

	"github.com/koykov/hash/xxhash"
	"github.com/koykov/pbtk/lsh"
)

func TestHash(t *testing.T) {
	h, err := NewHasher[[]byte](xxhash.Hasher64[[]byte]{})
	_ = err
	lsh.TestMe(t, h, lsh.TestDistHamming, 1.0)
}

func BenchmarkHash(b *testing.B) {
	h, err := NewHasher[[]byte](xxhash.Hasher64[[]byte]{})
	_ = err
	lsh.BenchmarkMe(b, h)
}
