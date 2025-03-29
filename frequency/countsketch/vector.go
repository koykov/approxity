package countsketch

import "io"

type vector interface {
	add(pos uint64, delta int64) error
	estimate(pos uint64) int64
	reset()
	readFrom(r io.Reader) (int64, error)
	writeTo(w io.Writer) (int64, error)
}

type vecbufh struct {
	p    uintptr
	l, c int
}
