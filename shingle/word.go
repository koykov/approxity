package shingle

import "github.com/koykov/byteseq"

type word[T byteseq.Q] struct {
	base[T]
	k uint
}

func NewWord[T byteseq.Q](k uint, cleanSet string) Shingler[T] {
	sh := &word[T]{
		base: base[T]{cset: cleanSet},
		k:    k,
	}
	sh.init()
	return sh
}

func (sh *word[T]) Shingle(s T) []T {
	bcap := 1
	if sh.k > 0 {
		bcap = len(s) / 3 / int(sh.k)
	}
	buf := make([]T, 0, bcap)
	return sh.AppendShingle(buf, s)
}

func (sh *word[T]) AppendShingle(dst []T, s T) []T {
	b := sh.clean(s, true)
	sc := byteseq.B2Q[T](b)
	if len(b) < 2 {
		dst = append(dst, sc)
		return dst
	}
	lo, hi := 0, sh.k
	_, _ = sh.spc[len(sh.spc)-1], sc[len(sc)-1]
	for i := uint64(0); i < uint64(len(sh.spc))-uint64(sh.k); i++ {
		dst = trimq(dst, sc[sh.spc[lo]:sh.spc[hi]])
		lo++
		hi++
	}
	return dst
}

func (sh *word[T]) Each(s T, fn func(T)) {
	b := sh.clean(s, true)
	sc := byteseq.B2Q[T](b)
	if len(b) < 2 {
		fn(sc)
		return
	}
	lo, hi := 0, sh.k
	_, _ = sh.spc[len(sh.spc)-1], sc[len(sc)-1]
	for i := uint64(0); i < uint64(len(sh.spc))-uint64(sh.k); i++ {
		fn(sc[sh.spc[lo]:sh.spc[hi]])
		lo++
		hi++
	}
}

func (sh *word[T]) Reset() {
	sh.base.reset()
}

func trimq[T byteseq.Q](dst []T, s T) []T {
	if s[0] == ' ' {
		s = s[1:]
	}
	if len(s) > 0 {
		dst = append(dst, s)
	}
	return dst
}
