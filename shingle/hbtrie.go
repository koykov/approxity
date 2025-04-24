package shingle

// Hierarchical Bitmask Trie implementation to check set contains a rune.
type hbtrie [16]uint64

func (t *hbtrie) set(r rune) {
	b0, b1, b2, b3 := uint8(r>>24), uint8(r>>16), uint8(r>>8), uint8(r)
	t[b0/64] |= 1 << (b0 % 64)
	t[4+b1/64] |= 1 << (b1 % 64)
	t[8+b2/64] |= 1 << (b2 % 64)
	t[12+b3/64] |= 1 << (b3 % 64)
}

func (t *hbtrie) contains(r rune) bool {
	b0, b1, b2, b3 := uint8(r>>24), uint8(r>>16), uint8(r>>8), uint8(r)
	if t[b0/64]&(1<<(b0%64)) == 0 {
		return false
	}
	if t[4+b1/64]&(1<<(b1%64)) == 0 {
		return false
	}
	if t[8+b2/64]&(1<<(b2%64)) == 0 {
		return false
	}
	if t[12+b3/64]&(1<<(b3%64)) == 0 {
		return false
	}
	return true
}

func (t *hbtrie) reset() {
	(*t)[0], (*t)[1], (*t)[2], (*t)[3] = 0, 0, 0, 0
}
