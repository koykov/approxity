package dlcsketch

import (
	"fmt"
	"testing"
)

func TestOptimal(t *testing.T) {
	stages := []struct {
		n, d uint64
		cp   float64
		m    uint64
	}{
		{100000, 2, 0.1, 131714},
		{500000, 3, 0.5, 240449},
		{1000000, 2, 0.1, 1317147},
		{2000000, 4, 0.3, 415291},
		{10000000, 4, 0.1, 3085736},
	}
	for i := 0; i < len(stages); i++ {
		st := &stages[i]
		t.Run(fmt.Sprintf("%d_%d_%f", st.n, st.d, st.cp), func(t *testing.T) {
			m := optimalM(st.n, st.d, st.cp)
			if m != st.m {
				t.Errorf("expected %d, got %d", st.m, m)
			}
		})
	}
}
