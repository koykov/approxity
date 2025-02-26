package bloom

import (
	"fmt"
	"testing"
)

func TestOptimal(t *testing.T) {
	t.Run("size", func(t *testing.T) {
		type stage struct {
			key string
			m   uint64
			fpp float64
			n   uint64
		}
		stages := []stage{
			{"1000_0.01", 1000, 0.01, 9586},
			{"100000_0.001", 100000, 0.001, 1437759},
			{"1000000_0.005", 1000000, 0.005, 11027754},
			{"10000000_0.01", 10000000, 0.01, 95850584},
		}
		for _, st := range stages {
			t.Run(st.key, func(t *testing.T) {
				n := OptimalSize(st.m, st.fpp)
				if n != st.n {
					t.Errorf("OptimalSize(%d, %f) = %d, want %d", st.m, st.fpp, n, st.n)
				}
			})
		}
	})
	t.Run("number hash functions", func(t *testing.T) {
		type stage struct {
			m, n, k uint64
		}
		var stages = []stage{
			{1000, 9586, 7},
			{100000, 1437759, 10},
			{1000000, 11027754, 8},
			{10000000, 95850584, 7},
		}
		for _, st := range stages {
			t.Run(fmt.Sprintf("%d_%d", st.m, st.n), func(t *testing.T) {
				k := OptimalNumberHashFunctions(st.m, st.n)
				if k != st.k {
					t.Errorf("OptimalNumberHashFunctions(%d, %d) = %d, want %d", st.m, st.n, k, st.k)
				}
			})
		}
	})
}
