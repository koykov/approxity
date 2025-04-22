package shingle

import (
	"unicode/utf8"

	"github.com/koykov/byteseq"
)

type char[T byteseq.Q] struct {
	base[T]
	k uint
	w []uint64
}

func NewChar[T byteseq.Q](k uint, cleanSet string) Shingler[T] {
	sh := &char[T]{
		base: base[T]{cset: cleanSet},
		k:    k,
	}
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
	b := sh.clean(s)
	sc := byteseq.B2Q[T](b)
	if uint(len(b)) <= sh.k {
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
	lo, hi := uint64(0), uint64(sh.k)
	_, _ = sh.w[len(sh.w)-1], sc[len(sc)-1]
	for i := uint64(0); i < uint64(len(sh.w))-uint64(sh.k); i++ {
		dst = append(dst, sc[sh.w[lo]:sh.w[hi]])
		lo++
		hi++
	}
	dst = append(dst, sc[sh.w[lo]:])
	return dst
}

func (sh *char[T]) Each(s T, fn func(T)) {
	b := sh.clean(s)
	sc := byteseq.B2Q[T](b)
	if uint(len(b)) <= sh.k || sh.k == 0 {
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
	lo, hi := uint64(0), uint64(sh.k)
	_, _ = sh.w[len(sh.w)-1], sc[len(sc)-1]
	for i := uint64(0); i < uint64(len(sh.w))-uint64(sh.k); i++ {
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
