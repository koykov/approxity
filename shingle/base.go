package shingle

import (
	"unicode/utf8"

	"github.com/koykov/byteseq"
)

type base[T byteseq.Q] struct {
	cset   string
	ctrie  hbtrie
	ctable map[rune]struct{}
	cbuf   []byte
	spc    []int
}

func (b *base[T]) init() {
	if b.ctable == nil {
		b.ctable = make(map[rune]struct{})
	}
	for _, c := range b.cset {
		b.ctable[c] = struct{}{}
		b.ctrie.set(c)
	}
}

func (b *base[T]) clean(s T, collapseSpaces bool) []byte {
	ss := byteseq.Q2S(s)
	b.spc = append(b.spc, 0)
	var space, pspace bool
	for i, c := range ss {
		if b.ctrie.contains(c) {
			// may be wrong, clarification required
			if _, ok := b.ctable[c]; ok {
				space = i > 0 && i < len(ss)-1 && ss[i-1] != ' ' && ss[i+1] != ' ' && ss[i] != '\''
				continue
			}
		}
		if space {
			b.cbuf = append(b.cbuf, ' ')
			b.spc = append(b.spc, spcp(len(b.cbuf)))
			space = false
		}
		if c == ' ' && pspace && collapseSpaces {
			continue
		}
		b.cbuf = utf8.AppendRune(b.cbuf, c)
		if pspace = c == ' '; pspace {
			b.spc = append(b.spc, spcp(len(b.cbuf)))
		}
	}
	b.spc = append(b.spc, len(b.cbuf))
	return b.cbuf
}

func (b *base[T]) reset() {
	b.cbuf = b.cbuf[:0]
	b.spc = b.spc[:0]
}

func spcp(p int) int {
	if p == 0 {
		return 0
	}
	return p - 1
}
