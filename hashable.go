package approxity

import (
	"bytes"
	"cmp"
	"slices"
	"strings"
)

// Hashable is a constraint that permits any type that can be hashed.
type Hashable interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
	~float32 | ~float64 |
	~string | ~[]byte | ~[]rune
}

func cmph[T Hashable](a, b T) int {
	switch x := any(a).(type) {
	case int:
		return cmp.Compare(x, any(b).(int))
	case int8:
		return cmp.Compare(x, any(b).(int8))
	case int16:
		return cmp.Compare(x, any(b).(int16))
	case int32:
		return cmp.Compare(x, any(b).(int32))
	case int64:
		return cmp.Compare(x, any(b).(int64))
	case uint:
		return cmp.Compare(x, any(b).(uint))
	case uint8:
		return cmp.Compare(x, any(b).(uint8))
	case uint16:
		return cmp.Compare(x, any(b).(uint16))
	case uint32:
		return cmp.Compare(x, any(b).(uint32))
	case uint64:
		return cmp.Compare(x, any(b).(uint64))
	case uintptr:
		return cmp.Compare(x, any(b).(uintptr))
	case float32:
		return cmp.Compare(x, any(b).(float32))
	case float64:
		return cmp.Compare(x, any(b).(float64))
	case string:
		return strings.Compare(x, any(b).(string))
	case []byte:
		return bytes.Compare(x, any(b).([]byte))
	case []rune:
		y := any(b).([]rune)
		n := min(len(x), len(y))
		for i := 0; i < n; i++ {
			if x[i] < y[i] {
				return -1
			}
			if x[i] > y[i] {
				return 1
			}
		}
		return 0
	}
	return 0
}

func Deduplicate[T Hashable](vals []T) []T {
	if len(vals) == 0 {
		return vals
	}
	slices.SortFunc(vals, cmph)
	var off int
	for i := 1; i < len(vals); i++ {
		if cmph(vals[i], vals[off]) != 0 {
			vals[off+1] = vals[i]
			off++
		}
	}
	return vals[:off+1]
}
