package cuckoo

import "unsafe"

const bucketsz = 4

type bucket uint32

func (b *bucket) add(fp byte) error {
	bb := b.b()
	for i := 0; i < bucketsz; i++ {
		if bb[i] == 0 {
			bb[i] = fp
			return nil
		}
	}
	return ErrFullBucket
}

func (b *bucket) set(i uint64, fp byte) error {
	bb := b.b()
	bb[i] = fp
	return nil
}

func (b *bucket) fpv(i uint64) byte {
	return b.b()[i]
}

func (b *bucket) fpi(fp byte) int {
	bb := b.b()
	for i := 0; i < bucketsz; i++ {
		if bb[i] == fp {
			return i
		}
	}
	return -1
}

func (b *bucket) b() *[bucketsz]byte {
	return (*[bucketsz]byte)(unsafe.Pointer(b))
}
