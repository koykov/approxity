package bloom

import (
	"encoding/binary"
	"io"
	"math"
	"unsafe"

	"github.com/koykov/approxity"
	"github.com/koykov/openrt"
)

const (
	cvectorDumpSignature = 0x5b0fc0b3cfae2b9
	cvectorDumpVersion   = 1.0
)

// Synchronous counting vector implementation.
type cvector struct {
	buf []uint32
	s   uint64
}

func (vec *cvector) Set(i uint64) bool {
	c := vec.buf[i/2]
	v0, v1 := uint16(c>>16), uint16(c)
	if i%2 == 0 {
		v0++
	} else {
		v1++
	}
	vec.buf[i/2] = uint32(v0)<<16 | uint32(v1)
	vec.s++
	return true
}

func (vec *cvector) Unset(i uint64) bool {
	c := vec.buf[i/2]
	v0, v1 := uint16(c>>16), uint16(c)
	if i%2 == 0 {
		v0 += math.MaxUint16
	} else {
		v1 += math.MaxUint16
	}
	vec.buf[i/2] = uint32(v0)<<16 | uint32(v1)
	vec.s += math.MaxUint16
	return true
}

func (vec *cvector) Get(i uint64) uint8 {
	c := vec.buf[i/2]
	v0, v1 := uint16(c>>16), uint16(c)
	var r bool
	if i%2 == 0 {
		r = v0 > 0
	} else {
		r = v1 > 0
	}
	if r {
		return 1
	}
	return 0
}

func (vec *cvector) Size() uint64 {
	return vec.s
}

func (vec *cvector) Capacity() uint64 {
	return uint64(len(vec.buf)) * 2
}

func (vec *cvector) Reset() {
	openrt.MemclrUnsafe(unsafe.Pointer(&vec.buf[0]), len(vec.buf)*4)
	vec.s = 0
}

func (vec *cvector) WriteTo(w io.Writer) (n int64, err error) {
	var (
		buf [32]byte
		m   int
	)
	binary.LittleEndian.PutUint64(buf[0:8], cvectorDumpSignature)
	binary.LittleEndian.PutUint64(buf[8:16], math.Float64bits(cvectorDumpVersion))
	binary.LittleEndian.PutUint64(buf[16:24], vec.s)
	if m, err = w.Write(buf[:]); err != nil {
		return int64(m), err
	}
	n += int64(m)

	var off int
	const blocksz = 4096
	var blk [blocksz]byte
	for i := 0; i < len(vec.buf); i++ {
		binary.LittleEndian.PutUint32(blk[off:], vec.buf[i])
		if off += 4; off == blocksz {
			m, err = w.Write(blk[:off])
			n += int64(m)
			if err != nil {
				return
			}
			if m < blocksz {
				err = io.ErrShortWrite
				return
			}
			off = 0
		}
	}
	if off > 0 {
		m, err = w.Write(blk[:off])
		n += int64(m)
	}
	return
}

func (vec *cvector) ReadFrom(r io.Reader) (n int64, err error) {
	var (
		buf [32]byte
		m   int
	)
	m, err = r.Read(buf[:])
	n += int64(m)
	if err != nil {
		return n, err
	}

	sign, ver, s := binary.LittleEndian.Uint64(buf[0:8]), binary.LittleEndian.Uint64(buf[8:16]),
		binary.LittleEndian.Uint64(buf[16:24])

	if sign != cvectorDumpSignature {
		return n, approxity.ErrInvalidSignature
	}
	if ver != math.Float64bits(cvectorDumpVersion) {
		return n, approxity.ErrVersionMismatch
	}
	vec.s = s
	vec.buf = vec.buf[:0]

	const blocksz = 4096
	var blk [blocksz]byte
	for {
		m, err = r.Read(blk[:])
		n += int64(m)
		if err != nil && err != io.EOF {
			return n, err
		}
		for i := 0; i < m; i += 4 {
			v := binary.LittleEndian.Uint32(blk[i:])
			vec.buf = append(vec.buf, v)
		}
		if err == io.EOF {
			err = nil
			break
		}
	}
	return
}

func newCvector(size uint64) *cvector {
	return &cvector{buf: make([]uint32, size/2+1)}
}
