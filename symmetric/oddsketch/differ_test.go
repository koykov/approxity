package oddsketch

import (
	"testing"

	"github.com/koykov/hash/xxhash"
	"github.com/koykov/pbtk/lsh/minhash"
	"github.com/koykov/pbtk/shingle"
	"github.com/koykov/pbtk/symmetric"
)

const (
	testSz  = 1e6
	testFPP = .01
)

var (
	testh       = xxhash.Hasher64[[]byte]{}
	testshc     = shingle.NewChar[[]byte](3, "") // 3-gram
	testshw     = shingle.NewWord[[]byte](2, "") // 2-word shingle
	testlshc, _ = minhash.NewHasher[[]byte](minhash.NewConfig[[]byte](testh, 50, testshc))
	testlshw, _ = minhash.NewHasher[[]byte](minhash.NewConfig[[]byte](testh, 50, testshw))
)

func TestDiffer(t *testing.T) {
	t.Run("char", func(t *testing.T) {
		d, err := NewDiffer[[]byte](NewConfig[[]byte](testSz, testFPP, testlshc))
		if err != nil {
			t.Fatal(err)
		}
		symmetric.TestMe(t, d, 0)
	})
	t.Run("word", func(t *testing.T) {
		d, err := NewDiffer[[]byte](NewConfig[[]byte](testSz, testFPP, testlshw))
		if err != nil {
			t.Fatal(err)
		}
		symmetric.TestMe(t, d, 0)
	})
}

func BenchmarkDiffer(b *testing.B) {
	b.Run("char", func(b *testing.B) {
		d, err := NewDiffer[[]byte](NewConfig[[]byte](testSz, testFPP, testlshc))
		if err != nil {
			b.Fatal(err)
		}
		symmetric.BenchMe(b, d)
	})
	b.Run("word", func(b *testing.B) {
		d, err := NewDiffer[[]byte](NewConfig[[]byte](testSz, testFPP, testlshw))
		if err != nil {
			b.Fatal(err)
		}
		symmetric.BenchMe(b, d)
	})
}
