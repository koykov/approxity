package pds

type Hasher interface {
	Sum64(data []byte) uint64
}
