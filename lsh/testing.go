package lsh

import (
	"fmt"
	"math/bits"
	"testing"

	"github.com/koykov/pbtk/simtest"
)

func TestMe[T []byte](t *testing.T, hash Hasher[T], distFn func([]uint64, []uint64, uint64) float64, numHashes uint64, expectAvgDist float64) {
	simtest.EachTestingDataset(func(_ int, ds *simtest.Dataset) {
		t.Run(ds.Name, func(t *testing.T) {
			var s, c float64
			for j := 0; j < len(ds.Tuples); j++ {
				tp := &ds.Tuples[j]

				hash.Reset()
				_ = hash.Add(tp.A)
				h0 := hash.Hash()

				hash.Reset()
				_ = hash.Add(tp.B)
				h1 := hash.Hash()

				dist := (64 - float64(distFn(h0, h1, numHashes))) / 64
				s += dist
				c++
			}
			if avg := s / c; avg > expectAvgDist {
				t.Errorf("avg dist = %f, expected %f", avg, expectAvgDist)
			}
		})
	})
}

func BenchMe[T []byte](b *testing.B, hash Hasher[T]) {
	stages := [][]byte{
		[]byte("foo"),
		[]byte("foobar"),
		[]byte("hello world"),
		[]byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit."),
		[]byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit. Mauris varius nisi erat, ac vulputate elit malesuada ut."),
		[]byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit. Mauris varius nisi erat, ac vulputate elit malesuada ut. Nulla facilisi. Vestibulum nec sapien nisl. Curabitur at elit fringilla, consectetur dui nec, maximus quam. Proin dui ipsum, venenatis nec est non, consectetur semper leo. Curabitur quis arcu ornare, malesuada nibh vel, maximus neque."),
	}
	for _, st := range stages {
		b.Run(fmt.Sprintf("add/%d", len(st)), func(b *testing.B) {
			b.SetBytes(int64(len(st)))
			b.ReportAllocs()
			b.ResetTimer()
			for j := 0; j < b.N; j++ {
				hash.Reset()
				_ = hash.Add(st)
			}
		})
	}
	for _, st := range stages {
		b.Run(fmt.Sprintf("hash/%d", len(st)), func(b *testing.B) {
			var buf []uint64
			b.SetBytes(int64(len(st)))
			b.ReportAllocs()
			b.ResetTimer()
			for j := 0; j < b.N; j++ {
				hash.Reset()
				_ = hash.Add(st)
				buf = hash.AppendHash(buf[:0])
				_ = buf
			}
		})
	}
}

func TestDistHamming(h0, h1 []uint64, _ uint64) (r float64) {
	n := max(len(h0), len(h1))
	for i := 0; i < n; i++ {
		var v0, v1 uint64
		if i < len(h0) {
			v0 = h0[i]
		}
		if i < len(h1) {
			v1 = h1[i]
		}
		r += float64(bits.OnesCount64(v0 ^ v1))
	}
	return r
}

func TestDistJaccard(h0, h1 []uint64, n uint64) (r float64) {
	if len(h1) < len(h0) {
		h0, h1 = h1, h0
	}
	for i := 0; i < len(h0); i++ {
		if h0[i] == h1[i] {
			r += 1
		}
	}
	if n == 0 {
		n = 1
	}
	return r / float64(n)
}
