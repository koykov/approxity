package quotient

import "testing"

func TestOptimal(t *testing.T) {
	tests := []struct {
		n       uint64
		fpp     float64
		m, q, r uint64
	}{
		{1000, 0.01, 2304, 11, 6},
		{1000, 0.001, 3072, 11, 9},
	}
	for _, test := range tests {
		m, q, r := optimalMQR(test.n, test.fpp, defaultLoadFactor)
		if m != test.m || q != test.q || r != test.r {
			t.Errorf("optimalMQR(%d, %f) = %d, %d, %d; want %d, %d, %d", test.n, test.fpp, m, q, r, test.m, test.q, test.r)
		}
	}
}
