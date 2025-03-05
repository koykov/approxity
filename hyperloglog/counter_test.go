package hyperloglog

import (
	"encoding/binary"
	"math"
	"testing"

	"github.com/koykov/hash/murmur/v3"
)

func TestCounter(t *testing.T) {
	const p = 18
	m := 1 << p
	count := 10 * m
	relative := 3.0 / math.Sqrt(float64(m))
	t.Run("count distincts", func(t *testing.T) {
		var buf [8]byte
		c, err := NewCounter(&Config{Precision: p, Hasher: murmur.Hasher64x64[[]byte]{Seed: 11400714819323198485}})
		if err != nil {
			t.Fatal(err)
		}
		for i := 0; i < 10; i++ {
			for j := uint64(1); j < uint64(count); j++ {
				binary.LittleEndian.PutUint64(buf[:], j)
				if err = c.Add(buf[:]); err != nil {
					t.Fatal(err)
				}
			}
		}
		e := c.Count()
		t.Log(e, e < uint64(float64(count)*(1+relative)) && e > uint64(float64(count)*(1-relative)))
	})
}
