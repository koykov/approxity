package cuckoo

import (
	"math"
	"sync/atomic"
)

// Concurrent ivector implementation.
type cnvector struct {
	buf []uint32
	lim uint64
	s   uint64
}

func (vec *cnvector) add(i uint64, fp byte) error {
	for k := uint64(0); k < vec.lim+1; k++ {
		for j := 0; j < bucketsz; j++ {
			if o := atomic.LoadUint32(&vec.buf[i]); o&vecmask[j] == 0 {
				n := o | uint32(fp)<<j
				if atomic.CompareAndSwapUint32(&vec.buf[i], o, n) {
					atomic.AddUint64(&vec.s, 1)
					return nil
				}
			}
		}
		return ErrFullBucket
	}
	return ErrWriteLimitReach
}

func (vec *cnvector) set(i, j uint64, fp byte) error {
	for k := uint64(0); k < vec.lim+1; k++ {
		o := atomic.LoadUint32(&vec.buf[i])
		n := o | uint32(fp)<<j
		if atomic.CompareAndSwapUint32(&vec.buf[i], o, n) {
			return nil
		}
	}
	return nil
}

func (vec *cnvector) unset(i uint64, fp byte) bool {
	for j := 0; j < bucketsz; j++ {
		if o := atomic.LoadUint32(&vec.buf[i]); o&vecmask[j] == uint32(fp)<<j {
			n := o & ^vecmask[j]
			if atomic.CompareAndSwapUint32(&vec.buf[i], o, n) {
				atomic.AddUint64(&vec.s, math.MaxUint64)
				return true
			}
		}
	}
	return false
}

func (vec *cnvector) fpv(i, j uint64) byte {
	return byte(atomic.LoadUint32(&vec.buf[i]) & vecmask[j] >> j)
}

func (vec *cnvector) fpi(i uint64, fp byte) int {
	for j := 0; j < bucketsz; j++ {
		if atomic.LoadUint32(&vec.buf[i])&vecmask[j] == uint32(fp)<<j {
			return j
		}
	}
	return -1
}

func (vec *cnvector) size() uint64 {
	return atomic.LoadUint64(&vec.s)
}

func (vec *cnvector) reset() {
	for i := 0; i < len(vec.buf); i++ {
		atomic.StoreUint32(&vec.buf[i], 0)
	}
}

func newCnvector(size, lim uint64) *cnvector {
	return &cnvector{
		buf: make([]uint32, size),
		lim: lim,
	}
}
