package shingle

import (
	"unicode/utf8"

	"github.com/koykov/byteseq"
)

type base[T byteseq.Q] struct {
	cset   string
	ctable map[rune]struct{}
	cbuf   []byte
}

func (b *base[T]) init() {
	if b.ctable == nil {
		b.ctable = make(map[rune]struct{})
	}
	for _, c := range b.cset {
		b.ctable[c] = struct{}{}
	}
}

func (b *base[T]) clean(s T) []byte {
	if len(b.cset) == 0 {
		b.cbuf = append(b.cbuf, s...)
		return b.cbuf
	}
	ss := byteseq.Q2S(s)
	var space bool
	for i, c := range ss {
		if _, ok := b.ctable[c]; ok {
			space = i > 0 && i < len(ss)-1 && ss[i-1] != ' ' && ss[i+1] != ' '
			continue
		}
		if space {
			b.cbuf = append(b.cbuf, ' ')
			space = false
		}
		b.cbuf = utf8.AppendRune(b.cbuf, c)
	}
	return b.cbuf
}

func (b *base[T]) reset() {
	b.cset = b.cset[:0]
	b.cbuf = b.cbuf[:0]
	for k := range b.ctable {
		delete(b.ctable, k)
	}
}
