package bloom

import "github.com/koykov/openrt"

type bitarray struct {
	buf []byte
}

func (b *bitarray) prealloc(size uint) {
	b.buf = make([]byte, size/8+1)
}

func (b *bitarray) set(i int) {
	b.buf[i/8] |= 1 << (i % 8)
}

func (b *bitarray) clear(i int) {
	b.buf[i/8] &^= 1 << (i % 8)
}

func (b *bitarray) get(i int) bool {
	return b.buf[i/8]&(1<<(i%8)) != 0
}

func (b *bitarray) reset() {
	openrt.Memclr(b.buf)
}
