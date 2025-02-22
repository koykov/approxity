package amq

import (
	"strconv"
	"unsafe"

	"github.com/koykov/x2bytes"
)

type Base struct{}

// HashSalt calculates hash sum of data + salt using given hasher.
func (b Base) HashSalt(hasher Hasher, data any, salt uint64) (_ uint64, err error) {
	return b.hash(hasher, data, salt, true)
}

// Hash calculates hash sum of data using given hasher.
func (b Base) Hash(hasher Hasher, data any) (_ uint64, err error) {
	return b.hash(hasher, data, 0, false)
}

func (b Base) hash(hasher Hasher, data any, salt uint64, saltext bool) (_ uint64, err error) {
	const bufsz = 64
	var a [bufsz]byte
	var h struct {
		ptr      uintptr
		len, cap int
	}
	h.ptr, h.cap = uintptr(unsafe.Pointer(&a)), bufsz
	buf := *(*[]byte)(unsafe.Pointer(&h))

	switch x := data.(type) {
	case []byte:
		buf = append(buf, x...)
	case *[]byte:
		buf = append(buf, *x...)
	case string:
		buf = append(buf, x...)
	case *string:
		buf = append(buf, *x...)
	default:
		if buf, err = x2bytes.ToBytes(buf, data); err != nil {
			return 0, err
		}
	}
	if saltext {
		buf = strconv.AppendUint(buf, salt, 10)
	}
	return hasher.Sum64(buf), err
}
