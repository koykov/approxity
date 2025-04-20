package shingle

import (
	"math"
	"unsafe"

	"github.com/koykov/byteseq"
	"github.com/koykov/openrt"
)

type base[T byteseq.Q] struct {
	cset   []byte
	ctable [math.MaxUint8]bool
	cbuf   []byte
}

func (b *base[T]) init() {
	for _, c := range b.cset {
		b.ctable[c] = true
	}
}

func (b *base[T]) clean(s T) []byte {
	if len(b.cset) == 0 {
		b.cbuf = append(b.cbuf, s...)
		return b.cbuf
	}
	_ = b.ctable[math.MaxUint8-1]
	for i := 0; i < len(s); i++ {
		if b.ctable[s[i]] {
			continue
		}
		b.cbuf = append(b.cbuf, s[i])
	}
	return b.cbuf
}

func (b *base[T]) reset() {
	b.cset = b.cset[:0]
	b.cbuf = b.cbuf[:0]
	openrt.MemclrUnsafe(unsafe.Pointer(&b.ctable[0]), len(b.ctable))
}
