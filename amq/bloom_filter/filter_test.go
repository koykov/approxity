package bloom

import (
	"os"
	"testing"

	"github.com/koykov/hash/xxhash"
	"github.com/koykov/pbtk/amq"
)

const (
	testSz  = 1e6
	testFPP = .01
)

var testh = xxhash.Hasher64[[]byte]{}

func TestFilter(t *testing.T) {
	t.Run("sync", func(t *testing.T) {
		f, err := NewFilter[[]byte](NewConfig(testSz, testFPP, testh))
		if err != nil {
			t.Fatal(err)
		}
		amq.TestMe(t, f)
	})
	t.Run("concurrent", func(t *testing.T) {
		f, err := NewFilter[[]byte](NewConfig(testSz, testFPP, testh).
			WithConcurrency().WithWriteAttemptsLimit(5))
		if err != nil {
			t.Fatal(err)
		}
		amq.TestMeConcurrently(t, f)
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
		t.Run("sync", func(t *testing.T) {
			f, _ := NewFilter[string](NewConfig(10, 0.01, testh))
			testWrite(t, f, "testdata/filter.bin", 45)
		})
		t.Run("concurrent", func(t *testing.T) {
			f, _ := NewFilter[string](NewConfig(10, 0.01, testh).WithConcurrency())
			testWrite(t, f, "testdata/concurrent_filter.bin", 56)
		})
	})
	t.Run("reader", func(t *testing.T) {
		testRead := func(t *testing.T, f amq.Filter[string], path string, expect int64) {
			fh, err := os.OpenFile(path, os.O_RDONLY, 0644)
			if err != nil {
				t.Fatal(err)
			}
			defer func() { _ = fh.Close() }()
			n, err := f.ReadFrom(fh)
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
		t.Run("sync", func(t *testing.T) {
			f, _ := NewFilter[string](NewConfig(10, 0.01, testh))
			testRead(t, f, "testdata/filter.bin", 45)
		})
		t.Run("concurrent", func(t *testing.T) {
			f, _ := NewFilter[string](NewConfig(10, 0.01, testh).WithConcurrency())
			testRead(t, f, "testdata/concurrent_filter.bin", 56)
		})
	})
}

func TestCountingFilter(t *testing.T) {
	t.Run("sync", func(t *testing.T) {
		f, err := NewCountingFilter[[]byte](NewConfig(testSz, testFPP, testh))
		if err != nil {
			t.Fatal(err)
		}
		amq.TestMe(t, f)
	})
	t.Run("concurrent", func(t *testing.T) {
		f, err := NewCountingFilter[[]byte](NewConfig(testSz, testFPP, testh).
			WithConcurrency().WithWriteAttemptsLimit(5))
		if err != nil {
			t.Fatal(err)
		}
		amq.TestMeConcurrently(t, f)
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
		t.Run("sync", func(t *testing.T) {
			f, _ := NewCountingFilter[string](NewConfig(10, 0.01, testh))
			testWrite(t, f, "testdata/counting_filter.bin", 228)
		})
		t.Run("concurrent", func(t *testing.T) {
			f, _ := NewCountingFilter[string](NewConfig(10, 0.01, testh).WithConcurrency())
			testWrite(t, f, "testdata/concurrent_counting_filter.bin", 236)
		})
	})
	t.Run("reader", func(t *testing.T) {
		testRead := func(t *testing.T, f amq.Filter[string], path string, expect int64) {
			fh, err := os.OpenFile(path, os.O_RDONLY, 0644)
			if err != nil {
				t.Fatal(err)
			}
			defer func() { _ = fh.Close() }()
			n, err := f.ReadFrom(fh)
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
		t.Run("sync", func(t *testing.T) {
			f, _ := NewCountingFilter[string](NewConfig(10, 0.01, testh))
			testRead(t, f, "testdata/counting_filter.bin", 228)
		})
		t.Run("concurrent", func(t *testing.T) {
			f, _ := NewCountingFilter[string](NewConfig(10, 0.01, testh).WithConcurrency())
			testRead(t, f, "testdata/concurrent_counting_filter.bin", 236)
		})
	})
}

func BenchmarkFilter(b *testing.B) {
	b.Run("sync", func(b *testing.B) {
		f, err := NewFilter[[]byte](NewConfig(testSz, testFPP, testh))
		if err != nil {
			b.Fatal(err)
		}
		amq.BenchMe(b, f)
	})
	b.Run("concurrent", func(b *testing.B) {
		f, err := NewFilter[[]byte](NewConfig(testSz, testFPP, testh).
			WithConcurrency().WithWriteAttemptsLimit(5))
		if err != nil {
			b.Fatal(err)
		}
		amq.BenchMeConcurrently(b, f)
	})
}

func BenchmarkCountingFilter(b *testing.B) {
	b.Run("sync", func(b *testing.B) {
		f, err := NewCountingFilter[[]byte](NewConfig(testSz, testFPP, testh))
		if err != nil {
			b.Fatal(err)
		}
		amq.BenchMe(b, f)
	})
	b.Run("concurrent", func(b *testing.B) {
		f, err := NewCountingFilter[[]byte](NewConfig(testSz, testFPP, testh).
			WithConcurrency().WithWriteAttemptsLimit(5))
		if err != nil {
			b.Fatal(err)
		}
		amq.BenchMeConcurrently(b, f)
	})
}
