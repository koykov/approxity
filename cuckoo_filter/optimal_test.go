package cuckoo

import (
	"testing"
)

func TestOptimal(t *testing.T) {
	t.Run("size", func(t *testing.T) {
		type stage struct {
			key string
			m   uint64
			n   uint64
		}
		stages := []stage{
			{"1000", 1000, 256},
			{"100000", 100000, 32768},
			{"1000000", 1000000, 262144},
			{"10000000", 10000000, 4194304},
		}
		for _, st := range stages {
			t.Run(st.key, func(t *testing.T) {
				n := optimalM(st.m)
				if n != st.n {
					t.Errorf("optimalM(%d) = %d, want %d", st.m, n, st.n)
				}
			})
		}
	})
}
