package bloom

import "github.com/koykov/openrt"

type bitarray struct {
	buf []byte
}

func (b *bitarray) prealloc(size uint) *bitarray {
	b.buf = make([]byte, size/8+1)
	return b
}

func (b *bitarray) set(i int) *bitarray {
	b.buf[i/8] |= 1 << (i % 8)
	return b
}

func (b *bitarray) clear(i int) *bitarray {
	b.buf[i/8] &^= 1 << (i % 8)
	return b
}

func (b *bitarray) get(i int) uint8 {
	return (b.buf[i/8] & (1 << (i % 8))) >> (i % 8)
}

func (b *bitarray) reset() *bitarray {
	openrt.Memclr(b.buf)
	return b
}
