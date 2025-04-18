package pbtk

import (
	"strconv"
	"unicode/utf8"
	"unsafe"
)

type Base[T Hashable] struct{}

// HashSalt calculates hash sum of data + salt using given hasher.
func (b Base[T]) HashSalt(hasher Hasher, data T, salt uint64) (uint64, error) {
	h, err := b.hash(hasher, nil, data, salt, true)
	return h[0], err
}

// Hash calculates hash sum of data using given hasher.
func (b Base[T]) Hash(hasher Hasher, data T) (uint64, error) {
	h, err := b.hash(hasher, nil, data, 0, false)
	return h[0], err
}

// HashSalt128 calculates 128-bit hash sum of data + salt using given hasher.
func (b Base[T]) HashSalt128(hasher Hasher128, data T, salt uint64) ([2]uint64, error) {
	return b.hash(nil, hasher, data, salt, true)
}

// Hash128 calculates 128-bit hash sum of data using given hasher.
func (b Base[T]) Hash128(hasher Hasher128, data T) ([2]uint64, error) {
	return b.hash(nil, hasher, data, 0, false)
}

func (b Base[T]) hash(hasher Hasher, hasher128 Hasher128, data T, salt uint64, saltext bool) (r [2]uint64, err error) {
	const bufsz = 64
	var a [bufsz]byte
	var h struct {
		ptr      uintptr
		len, cap int
	}
	h.ptr, h.cap = uintptr(unsafe.Pointer(&a)), bufsz
	buf := *(*[]byte)(unsafe.Pointer(&h))

	switch x := any(data).(type) {
	// int
	case int:
		buf = strconv.AppendInt(buf, int64(x), 10)
	case int8:
		buf = strconv.AppendInt(buf, int64(x), 10)
	case int16:
		buf = strconv.AppendInt(buf, int64(x), 10)
	case int32:
		buf = strconv.AppendInt(buf, int64(x), 10)
	case int64:
		buf = strconv.AppendInt(buf, x, 10)
	// uint
	case uint:
		buf = strconv.AppendUint(buf, uint64(x), 10)
	case uint8:
		buf = strconv.AppendUint(buf, uint64(x), 10)
	case uint16:
		buf = strconv.AppendUint(buf, uint64(x), 10)
	case uint32:
		buf = strconv.AppendUint(buf, uint64(x), 10)
	case uint64:
		buf = strconv.AppendUint(buf, x, 10)
	case uintptr:
		buf = strconv.AppendUint(buf, uint64(x), 10)
	// float
	case float32:
		buf = strconv.AppendFloat(buf, float64(x), 'f', -1, 32)
	case float64:
		buf = strconv.AppendFloat(buf, x, 'f', -1, 64)
	// byteseq
	case []byte:
		if len(x) > cap(buf) {
			buf = x
		} else {
			buf = append(buf, x...)
		}
	case string:
		buf = append(buf, x...)
	case []rune:
		if n := len(x); n > 0 {
			_ = x[n-1]
			for i := 0; i < n; i++ {
				buf = utf8.AppendRune(buf, x[i])
			}
		}
	default:
		return r, ErrEncoding
	}
	if saltext {
		buf = strconv.AppendUint(buf, salt, 10)
	}
	switch {
	case hasher != nil:
		r[0] = hasher.Sum64(buf)
		return r, err
	case hasher128 != nil:
		r = hasher128.Sum128(buf)
		return r, err
	default:
		return r, ErrEncoding
	}
}
