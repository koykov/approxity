package hyperloglog

import "io"

// dense vector interface
type vector interface {
	add(idx uint64, val uint8) error
	estimate() (float64, float64)
	capacity() uint64
	size() uint64
	reset()
	writeTo(w io.Writer) (n int64, err error)
	readFrom(r io.Reader) (n int64, err error)
}
