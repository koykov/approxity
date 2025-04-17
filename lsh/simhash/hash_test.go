package simhash

import (
	"testing"

	"github.com/koykov/hash/xxhash"
	"github.com/koykov/pbtk/lsh"
)

func TestHash(t *testing.T) {
	h, err := NewHasher[string](xxhash.Hasher64[[]byte]{})
	_ = err
	lsh.TestMe(t, h)
}
