package cuckoo

import (
	"os"
	"testing"

	"github.com/koykov/amq"
	"github.com/koykov/hash/xxhash"
)

const testsz = 1e7

var testh = xxhash.Hasher64[[]byte]{}

func TestFilter(t *testing.T) {
	t.Run("sync", func(t *testing.T) {
		f, err := NewFilter(NewConfig(testsz, testh))
		if err != nil {
			t.Fatal(err)
		}
		amq.TestFilter(t, f)
	})
	t.Run("concurrent", func(t *testing.T) {
		f, err := NewFilter(NewConfig(testsz, testh).
			WithConcurrency().WithWriteAttemptsLimit(5))
		if err != nil {
			t.Fatal(err)
		}
		amq.TestFilterConcurrently(t, f)
	})
	t.Run("writer", func(t *testing.T) {
		testWrite := func(t *testing.T, f amq.Filter, path string, expect int64) {
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
			f, _ := NewFilter(NewConfig(10, testh))
			testWrite(t, f, "testdata/filter.bin", 40)
		})
		t.Run("concurrent", func(t *testing.T) {
			f, _ := NewFilter(NewConfig(10, testh).WithConcurrency())
			testWrite(t, f, "testdata/concurrent_filter.bin", 40)
		})
	})
	t.Run("reader", func(t *testing.T) {
		testRead := func(t *testing.T, f amq.Filter, path string, expect int64) {
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
			f, _ := NewFilter(NewConfig(10, testh))
			testRead(t, f, "testdata/filter.bin", 40)
		})
		t.Run("concurrent", func(t *testing.T) {
			f, _ := NewFilter(NewConfig(10, testh).WithConcurrency())
			testRead(t, f, "testdata/concurrent_filter.bin", 40)
		})
	})
}

func BenchmarkFilter(b *testing.B) {
	b.Run("sync", func(b *testing.B) {
		f, err := NewFilter(NewConfig(testsz, testh))
		if err != nil {
			b.Fatal(err)
		}
		amq.BenchFilter(b, f)
	})
	b.Run("concurrent", func(b *testing.B) {
		f, err := NewFilter(NewConfig(testsz, testh).
			WithConcurrency().WithWriteAttemptsLimit(5))
		if err != nil {
			b.Fatal(err)
		}
		amq.BenchFilterConcurrently(b, f)
	})
}
