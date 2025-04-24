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

	m0 := uint64(1 << (b0 % 64))
	m1 := uint64(1 << (b1 % 64))
	m2 := uint64(1 << (b2 % 64))
	m3 := uint64(1 << (b3 % 64))

	_ = t[15]
	l0 := t[b0/64]
	l1 := t[4+b1/64]
	l2 := t[8+b2/64]
	l3 := t[12+b3/64]

	return (l0&m0 != 0) && (l1&m1 != 0) && (l2&m2 != 0) && (l3&m3 != 0)
}

func (t *hbtrie) reset() {
	for i := 0; i < 16; i++ {
		(*t)[i] = 0
	}
}
