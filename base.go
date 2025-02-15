package amq

import (
	"strconv"
	"unsafe"

	"github.com/koykov/x2bytes"
)

type Base struct{}

// Hash calculates hash sum of data + seed using given hasher.
func (Base) Hash(hasher Hasher, data any, seed uint64) (_ uint64, err error) {
	const bufsz = 128
	var a [bufsz]byte
	var h struct {
		ptr      uintptr
		len, cap int
	}
	h.ptr, h.cap = uintptr(unsafe.Pointer(&a)), bufsz
	buf := *(*[]byte)(unsafe.Pointer(&h))

	if buf, err = x2bytes.ToBytes(buf, data); err != nil {
		return 0, err
	}
	buf = strconv.AppendUint(buf, seed, 10)
	return hasher.Sum64(buf), err
}
