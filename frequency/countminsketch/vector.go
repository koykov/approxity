package countminsketch

import "io"

type vector interface {
	add(uint64) error
	estimate(uint64) uint64
	reset()
	readFrom(io.Reader) (int64, error)
	writeTo(io.Writer) (int64, error)
}
