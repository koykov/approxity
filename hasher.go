package approxity

type Hasher interface {
	Sum64(data []byte) uint64
}
