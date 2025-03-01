package cuckoo

import (
	"encoding/binary"
	"io"
	"math"
	"unsafe"

	"github.com/koykov/openrt"
)

const (
	vectorDumpSignature = 0x19329bb7706377b1
	vectorDumpVersion   = 1.0
)

type ivector interface {
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

// Synchronized ivector implementation.
type vector struct {
	buf []uint32
	s   uint64
}

func (vec *vector) add(i uint64, fp byte) error {
	for j := 0; j < bucketsz; j++ {
		if vec.buf[i]&vecmask[j] == 0 {
			vec.buf[i] |= uint32(fp) << (j * 8)
			vec.s++
			return nil
		}
	}
	return ErrFullBucket
}

func (vec *vector) set(i, j uint64, fp byte) error {
	vec.buf[i] |= uint32(fp) << (j * 8)
	return nil
}

func (vec *vector) unset(i uint64, fp byte) bool {
	for j := 0; j < bucketsz; j++ {
		if vec.buf[i]&vecmask[j] == uint32(fp)<<(j*8) {
			vec.buf[i] &= ^vecmask[j]
			vec.s--
			return true
		}
	}
	return false
}

func (vec *vector) fpv(i, j uint64) byte {
	return byte(vec.buf[i] & vecmask[j] >> (j * 8))
}

func (vec *vector) fpi(i uint64, fp byte) int {
	for j := 0; j < bucketsz; j++ {
		if vec.buf[i]&vecmask[j] == uint32(fp)<<(j*8) {
			return j
		}
	}
	return -1
}

func (vec *vector) capacity() uint64 {
	return uint64(len(vec.buf))
}

func (vec *vector) size() uint64 {
	return vec.s
}

func (vec *vector) reset() {
	openrt.MemclrUnsafe(unsafe.Pointer(&vec.buf[0]), len(vec.buf)*bucketsz)
}

func (vec *vector) writeTo(w io.Writer) (n int64, err error) {
	var (
		buf [24]byte
		m   int
	)
	binary.LittleEndian.PutUint64(buf[0:8], vectorDumpSignature)
	binary.LittleEndian.PutUint64(buf[8:16], math.Float64bits(vectorDumpVersion))
	binary.LittleEndian.PutUint64(buf[16:24], vec.s)
	m, err = w.Write(buf[:])
	n += int64(m)
	if err != nil {
		return int64(m), err
	}

	var h struct {
		p    uintptr
		l, c int
	}
	h.p = uintptr(unsafe.Pointer(&vec.buf[0]))
	h.l, h.c = len(vec.buf)*4, cap(vec.buf)*4
	payload := *(*[]byte)(unsafe.Pointer(&h))
	m, err = w.Write(payload)
	n += int64(m)
	return
}

func (vec *vector) readFrom(r io.Reader) (n int64, err error) {
	var (
		buf [24]byte
		m   int
	)
	m, err = r.Read(buf[:])
	n += int64(m)
	if err != nil {
		return n, err
	}

	sign, ver, s := binary.LittleEndian.Uint64(buf[0:8]), binary.LittleEndian.Uint64(buf[8:16]),
		binary.LittleEndian.Uint64(buf[16:24])

	if sign != vectorDumpSignature {
		return n, ErrInvalidSignature
	}
	if ver != math.Float64bits(vectorDumpVersion) {
		return n, ErrVersionMismatch
	}
	vec.s = s

	var h struct {
		p    uintptr
		l, c int
	}
	h.p = uintptr(unsafe.Pointer(&vec.buf[0]))
	h.l, h.c = len(vec.buf)*4, cap(vec.buf)*4
	payloadBuf := *(*[]byte)(unsafe.Pointer(&h))

	m, err = r.Read(payloadBuf)
	n += int64(m)
	if err == io.EOF {
		err = nil
	}
	return
}

func newVector(size uint64) *vector {
	return &vector{buf: make([]uint32, size)}
}
