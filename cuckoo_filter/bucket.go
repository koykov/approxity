package cuckoo

import "unsafe"

type bucket uint64

func (b *bucket) add(fp byte) error {
	bb := b.b()
	for i := 0; i < 8; i++ {
		if bb[i] == 0 {
			bb[i] = fp
			return nil
		}
	}
	return ErrFullBucket
}

func (b *bucket) set(i uint64, fp byte) error {
	b.b()[i] = fp
	return nil
}

func (b *bucket) get(i uint64) byte {
	return b.b()[i]
}

func (b *bucket) b() []byte {
	type sh struct {
		p    uintptr
		l, c int
	}
	h := sh{p: uintptr(unsafe.Pointer(b)), l: 8, c: 8}
	return *(*[]byte)(unsafe.Pointer(&h))
}
