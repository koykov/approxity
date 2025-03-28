package countsketch

import "testing"

func TestOptimal(t *testing.T) {
	stages := []struct {
		c, e float64
		w, d uint64
	}{
		{0.99, 0.01, 10000, 5},
		{0.99, 0.001, 1000000, 5},
		{0.99, 0.0001, 100000000, 5},
		{0.99, 0.00001, 10000000000, 5},
	}
	for _, stage := range stages {
		w, d := optimalWD(stage.c, stage.e)
		if w != stage.w || d != stage.d {
			t.Errorf("optimalWD(%f, %f) = %d, %d; want %d, %d", stage.c, stage.e, w, d, stage.w, stage.d)
		}
	}
}
