package hyperloglog

import (
	"math"
	"strconv"
	"testing"

	"github.com/koykov/hash/xxhash"
)

func TestCounter(t *testing.T) {
	const p = 18
	m := 1 << p
	count := 10 * m
	relative_error := 3.0 / math.Sqrt(float64(m))
	t.Run("count distincts", func(t *testing.T) {
		var buf []byte
		c, err := NewCounter(&Config{Precision: p, Hasher: xxhash.Hasher64[[]byte]{}})
		if err != nil {
			t.Fatal(err)
		}
		for i := 0; i < 10; i++ {
			for j := 0; j < count; j++ {
				buf = strconv.AppendInt(buf[:0], int64(j), 10)
				if err = c.Add(buf); err != nil {
					t.Fatal(err)
				}
			}
		}
		e := c.Count()
		t.Log(e < uint64(float64(count)*(1+relative_error)))
		t.Log(e > uint64(float64(count)*(1-relative_error)))
	})
}
