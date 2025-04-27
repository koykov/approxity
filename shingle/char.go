package shingle

import (
	"unicode/utf8"

	"github.com/koykov/byteseq"
)

type char[T byteseq.Q] struct {
	base[T]
	k uint64
	w []uint64
}

func NewChar[T byteseq.Q](k uint64, cleanSet string) Shingler[T] {
	sh := &char[T]{base: base[T]{cset: cleanSet}, k: k}
	sh.init()
	return sh
}

func (sh *char[T]) Shingle(s T) []T {
	bcap := 1
	if sh.k > 0 {
		bcap = len(s) / int(sh.k)
	}
	buf := make([]T, 0, bcap)
	return sh.AppendShingle(buf, s)
}

func (sh *char[T]) AppendShingle(dst []T, s T) []T {
	b := sh.clean(s, false)
	sc := byteseq.B2Q[T](b)
	if uint64(len(b)) <= sh.k || sh.k == 0 {
		dst = append(dst, sc)
		return dst
	}
	bl := uint64(len(b))
	_ = b[bl-1]
	for i := uint64(0); i < bl; {
		_, l := utf8.DecodeRune(b[i:])
		ul := uint64(l)
		sh.w = append(sh.w, i)
		i += ul
	}
	lo, hi := uint64(0), sh.k
	_, _ = sh.w[len(sh.w)-1], sc[len(sc)-1]
	for hi < uint64(len(sh.w)) {
		dst = append(dst, sc[sh.w[lo]:sh.w[hi]])
		lo++
		hi++
	}
	dst = append(dst, sc[sh.w[lo]:])
	return dst
}

func (sh *char[T]) Each(s T, fn func(T)) {
	b := sh.clean(s, false)
	sc := byteseq.B2Q[T](b)
	if uint64(len(b)) <= sh.k || sh.k == 0 {
		fn(sc)
		return
	}
	bl := uint64(len(b))
	_ = b[bl-1]
	for i := uint64(0); i < bl; {
		_, l := utf8.DecodeRune(b[i:])
		ul := uint64(l)
		sh.w = append(sh.w, i)
		i += ul
	}
	lo, hi := uint64(0), sh.k
	_, _ = sh.w[len(sh.w)-1], sc[len(sc)-1]
	for hi < uint64(len(sh.w)) {
		fn(sc[sh.w[lo]:sh.w[hi]])
		lo++
		hi++
	}
	fn(sc[sh.w[lo]:])
}

func (sh *char[T]) Reset() {
	sh.base.reset()
	sh.w = sh.w[:0]
}
