package cuckoo

import "unsafe"

type bucket uint64

func (b *bucket) add(fp byte) error {
	type sh struct {
		p    uintptr
		l, c int
	}
	h := sh{p: uintptr(unsafe.Pointer(b)), l: 8, c: 8}
	bb := *(*[]byte)(unsafe.Pointer(&h))
	for i := 0; i < 8; i++ {
		if bb[i] == 0 {
			bb[i] = fp
			return nil
		}
	}
	return ErrFullBucket
}
