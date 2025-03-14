package xor

import (
	"math"
	"os"
	"sync/atomic"
	"testing"

	"github.com/koykov/approxity"
	"github.com/koykov/approxity/amq"
	"github.com/koykov/hash/xxhash"
)

var testh = xxhash.Hasher64[[]byte]{}

func TestFilter(t *testing.T) {
	approxity.EachTestingDataset(func(i int, ds *approxity.TestingDataset[[]byte]) {
		t.Run(ds.Name, func(t *testing.T) {
			f, err := NewFilterWithKeys[[]byte](&Config{Hasher: testh}, ds.Positives)
			if err != nil {
				t.Fatal(err)
			}
			var falseNegative, falsePositive int
			for i := 0; i < len(ds.Negatives); i++ {
				if f.Contains(ds.Negatives[i]) {
					falsePositive++
				}
			}
			if falsePositive > 0 {
				// Just warn, it's OK to have small amount of false positives.
				t.Logf("%d of %d negatives (%d total) gives false positive value", falsePositive, len(ds.Negatives), len(ds.All))
			}
			for i := 0; i < len(ds.Positives); i++ {
				if !f.Contains(ds.Positives[i]) {
					falseNegative++
				}
			}
			if falseNegative > 0 {
				t.Errorf("%d of %d positives (%d total) gives false negative value", falseNegative, len(ds.Positives), len(ds.All))
			}
		})
	})
	t.Run("writer", func(t *testing.T) {
		testWrite := func(t *testing.T, f amq.Filter[string], path string, expect int64) {
			_ = f.Set("foobar")
			_ = f.Set("qwerty")
			fh, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				t.Fatal(err)
			}
			defer func() { _ = fh.Close() }()
			n, err := f.WriteTo(fh)
			if err != nil {
				t.Fatal(err)
			}
			if n != expect {
				t.Fatalf("expected %d bytes, got %d", expect, n)
			}
		}
		f, _ := NewFilterWithKeys[string](NewConfig(testh), []string{
			"foobar",
			"qwerty",
			"marquis",
			"warren",
		})
		testWrite(t, f, "testdata/filter.bin", 64)
	})
	t.Run("reader", func(t *testing.T) {
		testRead := func(t *testing.T, path string, expect int64) {
			fh, err := os.OpenFile(path, os.O_RDONLY, 0644)
			if err != nil {
				t.Fatal(err)
			}
			defer func() { _ = fh.Close() }()

			f, n, _ := NewFilterFromReader[string](NewConfig(testh), fh)
			if err != nil {
				t.Fatal(err)
			}
			if n != expect {
				t.Fatalf("expected %d bytes, got %d", expect, n)
			}
			if !f.Contains("foobar") || !f.Contains("qwerty") {
				t.Fatal("filter does not contain expected values")
			}
		}
		testRead(t, "testdata/filter.bin", 64)
	})
}

func BenchmarkFilter(b *testing.B) {
	b.Run("sync", func(b *testing.B) {
		approxity.EachTestingDataset(func(i int, ds *approxity.TestingDataset[[]byte]) {
			b.Run(ds.Name, func(b *testing.B) {
				f, _ := NewFilterWithKeys[[]byte](&Config{Hasher: testh}, ds.Positives)
				b.ReportAllocs()
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					f.Contains(ds.All[i%len(ds.All)])
				}
			})
		})
	})
	b.Run("concurrent", func(b *testing.B) {
		approxity.EachTestingDataset(func(i int, ds *approxity.TestingDataset[[]byte]) {
			b.Run(ds.Name, func(b *testing.B) {
				f, _ := NewFilterWithKeys[[]byte](&Config{Hasher: testh}, ds.Positives)
				b.ReportAllocs()
				b.RunParallel(func(pb *testing.PB) {
					var i uint64 = math.MaxUint64
					for pb.Next() {
						ci := atomic.AddUint64(&i, 1)
						f.Contains(ds.All[ci%uint64(len(ds.All))])
					}
				})
			})
		})
	})
}
