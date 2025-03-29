package countsketch

import (
	"os"
	"testing"

	"github.com/koykov/approxity/frequency"
	"github.com/koykov/hash/xxhash"
)

const (
	testConfidence = 0.999
	testEpsilon    = 0.001
)

var testh = xxhash.Hasher64[[]byte]{}

func TestEstimator(t *testing.T) {
	t.Run("sync", func(t *testing.T) {
		est, err := NewEstimator[[]byte](NewConfig(testConfidence, testEpsilon, testh))
		if err != nil {
			t.Fatal(err)
		}
		frequency.TestMe(t, frequency.NewTestSignedAdapter(est))
	})
	t.Run("concurrent", func(t *testing.T) {
		est, err := NewEstimator[[]byte](NewConfig(testConfidence, testEpsilon, testh).
			WithConcurrency())
		if err != nil {
			t.Fatal(err)
		}
		frequency.TestMeConcurrently(t, frequency.NewTestSignedAdapter(est))
	})
	t.Run("writer", func(t *testing.T) {
		testWrite := func(t *testing.T, est frequency.SignedEstimator[string], path string, expect int64) {
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
			t.Run("32", func(t *testing.T) {
				e, _ := NewEstimator[string](NewConfig(.99, .01, testh).
					WithCompact())
				testWrite(t, e, "testdata/estimator32.bin", 200016)
			})
			t.Run("64", func(t *testing.T) {
				e, _ := NewEstimator[string](NewConfig(.99, .01, testh))
				testWrite(t, e, "testdata/estimator64.bin", 400016)
			})
		})
		t.Run("concurrent", func(t *testing.T) {
			t.Run("32", func(t *testing.T) {
				e, _ := NewEstimator[string](NewConfig(.99, .01, testh).
					WithConcurrency().
					WithCompact())
				testWrite(t, e, "testdata/concurrent_estimator32.bin", 200016)
			})
			t.Run("64", func(t *testing.T) {
				e, _ := NewEstimator[string](NewConfig(.99, .01, testh).
					WithConcurrency())
				testWrite(t, e, "testdata/concurrent_estimator64.bin", 400016)
			})
		})
	})
	t.Run("reader", func(t *testing.T) {
		testRead := func(t *testing.T, est frequency.SignedEstimator[string], path string, expectBytes int64, expect0, expect1 int64) {
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
			e0, e1 := est.Estimate("foobar"), est.Estimate("qwerty")
			if e0 != expect0 {
				t.Errorf("expected %d estimate, got %d", expect0, e0)
			}
			if e1 != expect1 {
				t.Errorf("expected %d estimate, got %d", expect1, e1)
			}
		}
		t.Run("sync", func(t *testing.T) {
			t.Run("32", func(t *testing.T) {
				f, _ := NewEstimator[string](NewConfig(.99, .01, testh).
					WithCompact())
				testRead(t, f, "testdata/estimator32.bin", 200016, 1, 100)
			})
			t.Run("64", func(t *testing.T) {
				f, _ := NewEstimator[string](NewConfig(.99, .01, testh))
				testRead(t, f, "testdata/estimator64.bin", 400016, 1, 100)
			})
		})
		t.Run("concurrent", func(t *testing.T) {
			t.Run("32", func(t *testing.T) {
				f, _ := NewEstimator[string](NewConfig(.99, .01, testh).
					WithConcurrency().
					WithCompact())
				testRead(t, f, "testdata/concurrent_estimator32.bin", 200016, 1, 100)
			})
			t.Run("64", func(t *testing.T) {
				f, _ := NewEstimator[string](NewConfig(.99, .01, testh).
					WithConcurrency())
				testRead(t, f, "testdata/concurrent_estimator64.bin", 400016, 1, 100)
			})
		})
	})
}

func BenchmarkEstimator(b *testing.B) {
	b.Run("sync", func(b *testing.B) {
		est, err := NewEstimator[[]byte](NewConfig(testConfidence, testEpsilon, testh))
		if err != nil {
			b.Fatal(err)
		}
		frequency.BenchMe(b, frequency.NewTestSignedAdapter(est))
	})
	b.Run("concurrent", func(b *testing.B) {
		est, err := NewEstimator[[]byte](NewConfig(testConfidence, testEpsilon, testh).
			WithConcurrency().WithWriteAttemptsLimit(5))
		if err != nil {
			b.Fatal(err)
		}
		frequency.BenchMeConcurrently(b, frequency.NewTestSignedAdapter(est))
	})
}
