package bloom

type Hasher interface {
	Sum64(data string) uint64
}
