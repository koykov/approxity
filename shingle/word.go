package shingle

import "github.com/koykov/byteseq"

type word[T byteseq.Q] struct {
	base[T]
}

func NewWord[T byteseq.Q](cleanSet string) Shingler[T] {
	sh := &word[T]{base: base[T]{cset: cleanSet}}
	sh.init()
	return sh
}

func (sh *word[T]) Shingle(s T, k uint) []T {
	bcap := 1
	if k > 0 {
		bcap = len(s) / 3 / int(k)
	}
	buf := make([]T, 0, bcap)
	return sh.AppendShingle(buf, s, k)
}

func (sh *word[T]) AppendShingle(dst []T, s T, k uint) []T {
	b := sh.clean(s, true)
	sc := byteseq.B2Q[T](b)
	if len(b) < 2 {
		dst = append(dst, sc)
		return dst
	}
	lo, hi := 0, k
	_, _ = sh.spc[len(sh.spc)-1], sc[len(sc)-1]
	for hi < uint(len(sh.spc)) {
		dst = trimq(dst, sc[sh.spc[lo]:sh.spc[hi]])
		lo++
		hi++
	}
	return dst
}

func (sh *word[T]) Each(s T, k uint, fn func(T)) {
	b := sh.clean(s, true)
	sc := byteseq.B2Q[T](b)
	if len(b) < 2 {
		fn(sc)
		return
	}
	lo, hi := 0, k
	_, _ = sh.spc[len(sh.spc)-1], sc[len(sc)-1]
	for hi < uint(len(sh.spc)) {
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
