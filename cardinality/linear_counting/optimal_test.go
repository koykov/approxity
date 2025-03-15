package linear

import "testing"

func TestOptimal(t *testing.T) {
	var stages = []struct {
		n  uint64
		cp float64
		m  uint64
	}{
		{100, 0.01, 9950},             // 1.21 KB
		{1000, 0.01, 99500},           // 12.14 KB
		{10000, 0.01, 994992},         // 121.45 KB
		{100000, 0.01, 9949917},       // 1.18 MB
		{1000000, 0.01, 99499163},     // 11.86 MB
		{10000000, 0.01, 994991625},   // 118.61 MB
		{100000000, 0.01, 9949916248}, // 1.15 GB
	}
	for _, stage := range stages {
		m := optimalM(stage.n, stage.cp)
		if m != stage.m {
			t.Errorf("optimalM(%d, %f) = %d, want %d", stage.n, stage.cp, m, stage.m)
		}
	}
}
