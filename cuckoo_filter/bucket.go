package cuckoo

// 7 byte payload / 1 byte header
type bucket uint64

func (b *bucket) plen() uint64 {
	return uint64(*b << 7)
}

func (b *bucket) payload() uint64 {
	return uint64(*b >> 1)
}

func (b *bucket) add(fp byte) error {
	return nil
}
