package cmsketch

import (
	"context"
	"io"

	"github.com/koykov/bitset"
)

const (
	flagConservativeUpdate = iota
	flagLFU                // reserved
)

type vector interface {
	add(hkey, delta uint64) error
	estimate(hkey uint64) uint64
	decay(ctx context.Context, factor float64) error
	reset()
	readFrom(r io.Reader) (int64, error)
	writeTo(w io.Writer) (int64, error)
}

type basevec struct {
	w, d  uint64
	flags bitset.Bitset64
}

func vecpos(lo, hi uint32, w, i uint64) uint64 {
	return i*w + uint64(lo+hi*uint32(i))%w
}

type vecbufh struct {
	p    uintptr
	l, c int
}
