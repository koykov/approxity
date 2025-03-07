package cuckoo

var mask64 [65]uint64

func init() {
	for i := 0; i < 65; i++ {
		mask64[i] = (1 << i) - 1
	}
}
