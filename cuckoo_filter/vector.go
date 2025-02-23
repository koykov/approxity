package cuckoo

import (
	"math"
	"unsafe"

	"github.com/koykov/openrt"
)

type ivector interface {
	add(i uint64, fp byte) error
	set(i, j uint64, fp byte) error
	unset(i uint64, fp byte) bool
	fpv(i, j uint64) byte
	fpi(i uint64, fp byte) int
	size() uint64
	reset()
}

var vecmask = [bucketsz]uint32{
	math.MaxUint8,
	math.MaxUint8 << 8,
	math.MaxUint8 << 16,
	math.MaxUint8 << 24,
}

// Synchronized ivector implementation.
type vector struct {
	buf []uint32
	s   uint64
}

func (vec *vector) add(i uint64, fp byte) error {
	for j := 0; j < bucketsz; j++ {
		if vec.buf[i]&vecmask[j] == 0 {
			vec.buf[i] |= uint32(fp) << j
			vec.s++
			return nil
		}
	}
	return ErrFullBucket
}

func (vec *vector) set(i, j uint64, fp byte) error {
	vec.buf[i] |= uint32(fp << j)
	return nil
}

func (vec *vector) unset(i uint64, fp byte) bool {
	for j := 0; j < bucketsz; j++ {
		if vec.buf[i]&vecmask[j] == uint32(fp)<<j {
			vec.buf[i] &= ^vecmask[j]
			vec.s--
			return true
		}
	}
	return false
}

func (vec *vector) fpv(i, j uint64) byte {
	return byte(vec.buf[i] & vecmask[j] >> j)
}

func (vec *vector) fpi(i uint64, fp byte) int {
	for j := 0; j < bucketsz; j++ {
		if vec.buf[i]&vecmask[j] == uint32(fp)<<j {
			return j
		}
	}
	return -1
}

func (vec *vector) size() uint64 {
	return vec.s
}

func (vec *vector) reset() {
	openrt.MemclrUnsafe(unsafe.Pointer(&vec.buf[0]), len(vec.buf)*bucketsz)
}

func newVector(size uint64) *vector {
	return &vector{buf: make([]uint32, size)}
}
