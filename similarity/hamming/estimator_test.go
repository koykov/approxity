package hamming

import (
	"testing"

	"github.com/koykov/hash/xxhash"
	"github.com/koykov/pbtk/lsh/simhash"
	"github.com/koykov/pbtk/shingle"
	"github.com/koykov/pbtk/similarity"
)

var (
	testh       = xxhash.Hasher64[[]byte]{}
	testshc     = shingle.NewChar[[]byte](3, "") // 3-gram
	testshw     = shingle.NewWord[[]byte](2, "") // 2-word shingle
	testlshc, _ = simhash.NewHasher[[]byte](simhash.NewConfig[[]byte](testh, testshc))
	testlshw, _ = simhash.NewHasher[[]byte](simhash.NewConfig[[]byte](testh, testshw))
)

func TestEstimator(t *testing.T) {
	t.Run("char", func(t *testing.T) {
		e, err := NewEstimator[[]byte](NewConfig[[]byte](testlshc))
		if err != nil {
			t.Fatal(err)
		}
		similarity.TestMe(t, e, 1)
	})
	t.Run("word", func(t *testing.T) {
		e, err := NewEstimator[[]byte](NewConfig[[]byte](testlshw))
		if err != nil {
			t.Fatal(err)
		}
		similarity.TestMe(t, e, 1)
	})
}

func BenchmarkEstimator(b *testing.B) {
	b.Run("char", func(b *testing.B) {
		e, err := NewEstimator[[]byte](NewConfig[[]byte](testlshc))
		if err != nil {
			b.Fatal(err)
		}
		similarity.BenchMe(b, e)
	})
	b.Run("word", func(b *testing.B) {
		e, err := NewEstimator[[]byte](NewConfig[[]byte](testlshw))
		if err != nil {
			b.Fatal(err)
		}
		similarity.BenchMe(b, e)
	})
}
