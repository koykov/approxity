package approxity

type Hasher interface {
	Sum64(data []byte) uint64
}

type Hasher128 interface {
	Sum128(data []byte) [2]uint64
}
