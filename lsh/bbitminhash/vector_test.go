package bbitminhash

import (
	"math"
	"reflect"
	"testing"
)

var vecTestVals = []uint64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}

func TestVector(t *testing.T) {
	vec := newVector(5)

	vec.Grow(21)
	if vec.Len() != 21 {
		t.Errorf("expected 20, got %d", vec.Len())
	}

	vec.Memset(12)
	for i := 0; i < len(vecTestVals); i++ {
		if vec.Get(uint64(i)) != 12 {
			t.Errorf("pos %d expected 12, got %d", i, vec.Get(uint64(i)))
		}
	}

	vec.Reset()
	for i := 0; i < len(vecTestVals); i++ {
		if vec.Get(uint64(i)) != 0 {
			t.Errorf("expected 0, got %d", vec.Get(uint64(i)))
		}
	}

	for _, v := range vecTestVals {
		vec.Add(v)
	}
	var x []uint64
	x = vec.AppendAll(x)
	if !reflect.DeepEqual(x, vecTestVals) {
		t.Errorf("expected %v, got %v", vecTestVals, x)
	}
}

func BenchmarkVector(b *testing.B) {
	b.Run("add", func(b *testing.B) {
		vec := newVector(5)
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			vec.Reset()
			for _, v := range vecTestVals {
				vec.Add(v)
			}
		}
	})
	b.Run("set", func(b *testing.B) {
		vec := newVector(5)
		vec.Grow(10)
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			vec.SetMin(3, math.MaxUint32)
		}
	})
	b.Run("get", func(b *testing.B) {
		vec := newVector(5)
		vec.Memset(100)
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			vec.Get(3)
		}
	})
	b.Run("append all", func(b *testing.B) {
		vec := newVector(5)
		for _, v := range vecTestVals {
			vec.Add(v)
		}
		var x []uint64
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			x = vec.AppendAll(x[:0])
			_ = x
		}
	})
	b.Run("memset", func(b *testing.B) {
		vec := newVector(5)
		vec.Grow(10)
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			vec.Memset(math.MaxUint64)
		}
	})
}
