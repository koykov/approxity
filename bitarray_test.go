package bloom

import "testing"

func TestBitarray(t *testing.T) {
	prepare := func(size uint) *bitarray {
		var a bitarray
		a.prealloc(size).
			set(3).
			set(5).
			set(7).
			set(9)
		return &a
	}
	t.Run("set", func(t *testing.T) {
		a := prepare(10)
		if a.buf[0] != 168 || a.buf[1] != 2 {
			t.Fail()
		}
	})
	t.Run("get", func(t *testing.T) {
		a := prepare(10)
		chk := map[int]uint8{3: 1, 5: 1, 7: 1, 9: 1}
		for i := 0; i < 10; i++ {
			if chk[i] != a.get(i) {
				t.Fail()
			}
		}
	})
}

func BenchmarkBitarray(b *testing.B) {
	b.Run("set", func(b *testing.B) {
		var a bitarray
		a.prealloc(10)
		for i := 0; i < b.N; i++ {
			a.set(9)
		}
	})
	b.Run("get", func(b *testing.B) {
		var a bitarray
		a.prealloc(10)
		a.set(5)
		for i := 0; i < b.N; i++ {
			a.get(5)
		}
	})
}
