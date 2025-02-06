package bloom

type Hasher interface {
	Hash(data any) uint64
}
