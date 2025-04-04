package tinylfu

func encode(tdelta uint32, n uint32) uint64 {
	return (uint64(tdelta) << 32) | uint64(n)
}

func decode(val uint64) (tdelta uint32, n uint32) {
	tdelta = uint32(val >> 32)
	n = uint32(val & 0xFFFFFFFF)
	return
}
