package bbitminhash

import "testing"

var bbvTestVals = []uint64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}

func TestBbitvec(t *testing.T) {
	vec := newBbitvec(5)
	for _, v := range bbvTestVals {
		vec.add(v)
	}
	var c uint64
	vec.each(func(v uint64) {
		if v != c {
			t.Errorf("expected %d, got %d", c, v)
		}
		c++
	})
}

func BenchmarkBbitvec(b *testing.B) {
	b.Run("add", func(b *testing.B) {
		vec := newBbitvec(5)
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			vec.reset()
			for _, v := range bbvTestVals {
				vec.add(v)
			}
		}
	})
	b.Run("each", func(b *testing.B) {
		vec := newBbitvec(5)
		for _, v := range bbvTestVals {
			vec.add(v)
		}
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			vec.each(func(v uint64) {})
		}
	})
}
