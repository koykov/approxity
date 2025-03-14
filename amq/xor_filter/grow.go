package xor

func growu8(dst []uint8, ln uint64) []uint8 {
	if ln <= uint64(cap(dst)) {
		return dst[:ln]
	}
	dst = make([]uint8, ln)
	return dst
}

func growu32(dst []uint32, ln uint64) []uint32 {
	if ln <= uint64(cap(dst)) {
		return dst[:ln]
	}
	dst = make([]uint32, ln)
	return dst
}

func growu64(dst []uint64, ln uint64) []uint64 {
	if ln <= uint64(cap(dst)) {
		return dst[:ln]
	}
	dst = make([]uint64, ln)
	return dst
}
