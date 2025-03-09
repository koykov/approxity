package hyperloglog

import (
	"os"
	"testing"

	"github.com/koykov/approxity/cardinality"
	"github.com/koykov/hash/xxhash"
)

const testP = 18

var testh = xxhash.Hasher64[[]byte]{}

func TestEstimator(t *testing.T) {
	t.Run("sync", func(t *testing.T) {
		est, err := NewEstimator[[]byte](NewConfig(testP, testh))
		if err != nil {
			t.Fatal(err)
		}
		cardinality.TestMe(t, est, 0.005)
	})
	t.Run("concurrent", func(t *testing.T) {
		est, err := NewEstimator[[]byte](NewConfig(testP, testh).
			WithConcurrency().WithWriteAttemptsLimit(5))
		if err != nil {
			t.Fatal(err)
		}
		cardinality.TestMeConcurrently(t, est, 0.005)
	})
	t.Run("writer", func(t *testing.T) {
		testWrite := func(t *testing.T, est cardinality.Estimator[string], path string, expect int64) {
			_ = est.Add("foobar")
			for i := 0; i < 100; i++ {
				_ = est.Add("qwerty")
			}
			fh, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				t.Fatal(err)
			}
			defer func() { _ = fh.Close() }()
			n, err := est.WriteTo(fh)
			if err != nil {
				t.Fatal(err)
			}
			if n != expect {
				t.Fatalf("expected %d bytes, got %d", expect, n)
			}
		}
		t.Run("sync", func(t *testing.T) {
			f, _ := NewEstimator[string](NewConfig(18, testh))
			testWrite(t, f, "testdata/estimator.bin", 262184)
		})
		t.Run("concurrent", func(t *testing.T) {
			f, _ := NewEstimator[string](NewConfig(18, testh).WithConcurrency())
			testWrite(t, f, "testdata/concurrent_estimator.bin", 262184)
		})
	})
	t.Run("reader", func(t *testing.T) {
		testRead := func(t *testing.T, est cardinality.Estimator[string], path string, expectBytes int64, expectEst uint64) {
			fh, err := os.OpenFile(path, os.O_RDONLY, 0644)
			if err != nil {
				t.Fatal(err)
			}
			defer func() { _ = fh.Close() }()
			n, err := est.ReadFrom(fh)
			if err != nil {
				t.Fatal(err)
			}
			if n != expectBytes {
				t.Fatalf("expected %d bytes, got %d", expectBytes, n)
			}
			if e := est.Estimate(); e != expectEst {
				t.Errorf("expected %d estimate, got %d", expectEst, e)
			}
		}
		t.Run("sync", func(t *testing.T) {
			f, _ := NewEstimator[string](NewConfig(18, testh))
			testRead(t, f, "testdata/estimator.bin", 262184, 2)
		})
		t.Run("concurrent", func(t *testing.T) {
			f, _ := NewEstimator[string](NewConfig(18, testh).WithConcurrency())
			testRead(t, f, "testdata/concurrent_estimator.bin", 262184, 2)
		})
	})
}

func BenchmarkEstimator(b *testing.B) {
	b.Run("sync", func(b *testing.B) {
		est, err := NewEstimator[[]byte](NewConfig(testP, testh))
		if err != nil {
			b.Fatal(err)
		}
		cardinality.BenchMe(b, est)
	})
	b.Run("concurrent", func(b *testing.B) {
		est, err := NewEstimator[[]byte](NewConfig(testP, testh).
			WithConcurrency().WithWriteAttemptsLimit(5))
		if err != nil {
			b.Fatal(err)
		}
		cardinality.BenchMeConcurrently(b, est)
	})
}
