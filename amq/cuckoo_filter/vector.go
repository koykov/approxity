package cuckoo

import (
	"io"
	"math"
)

const bucketsz = 4

// Vector of uint32 buckets.
type vector interface {
	add(i uint64, fp byte) error
	set(i, j uint64, fp byte) error
	unset(i uint64, fp byte) bool
	fpv(i, j uint64) byte
	fpi(i uint64, fp byte) int
	capacity() uint64
	size() uint64
	reset()
	writeTo(w io.Writer) (n int64, err error)
	readFrom(r io.Reader) (n int64, err error)
}

var vecmask = [bucketsz]uint32{
	math.MaxUint8,
	math.MaxUint8 << 8,
	math.MaxUint8 << 16,
	math.MaxUint8 << 24,
}
