package cmsketch

import (
	"context"
	"encoding/binary"
	"io"
	"math"
	"sync/atomic"

	"github.com/koykov/bitset"
	"github.com/koykov/pbtk"
)

const (
	dumpSignatureConcurrent32 = 0x2CC21173A9E62A9D
	dumpVersionConcurrent32   = 1.0
)

// 32-bit version of concurrent vector implementation.
type cnvector32 struct {
	basevec
	lim uint64
	buf []uint32
}

func (vec *cnvector32) add(hkey, delta uint64) error {
	lo, hi := uint32(hkey>>32), uint32(hkey)

	switch {
	case vec.flags.CheckBit(flagConservativeUpdate):
		return vec.addCU(lo, hi, delta)
	case vec.flags.CheckBit(flagLFU):
		return vec.addLFU(lo, hi, delta)
	default:
		return vec.addClassic(lo, hi, delta)
	}
}

func (vec *cnvector32) addClassic(lo, hi uint32, delta uint64) error {
	for i := uint64(0); i < vec.d; i++ {
		pos := i*vec.w + uint64(lo+hi*uint32(i))%vec.w
		var ok bool
		var j uint64
		for j = 0; j < vec.lim+1; j++ {
			o := atomic.LoadUint32(&vec.buf[pos])
			n := o + uint32(delta)
			if ok = atomic.CompareAndSwapUint32(&vec.buf[pos], o, n); ok {
				break
			}
		}
		if !ok {
			return pbtk.ErrWriteLimitExceed
		}
	}
	return nil
}

func (vec *cnvector32) addCU(lo, hi uint32, delta uint64) error {
	var mn uint32 = math.MaxUint32
	for i := uint64(0); i < vec.d; i++ {
		mn = min(mn, atomic.LoadUint32(&vec.buf[vecpos(lo, hi, vec.w, i)]))
	}
	for i := uint64(0); i < vec.d; i++ {
		pos := vecpos(lo, hi, vec.w, i)
		if atomic.LoadUint32(&vec.buf[pos]) == mn {
			var ok bool
			var j uint64
			for j = 0; j < vec.lim+1; j++ {
				o := atomic.LoadUint32(&vec.buf[pos])
				n := o + uint32(delta)
				if ok = atomic.CompareAndSwapUint32(&vec.buf[pos], o, n); ok {
					break
				}
			}
			if !ok {
				return pbtk.ErrWriteLimitExceed
			}
		}
	}
	return nil
}

func (vec *cnvector32) addLFU(lo, hi uint32, delta uint64) error {
	return vec.addClassic(lo, hi, delta)
}

func (vec *cnvector32) estimate(hkey uint64) (r uint64) {
	lo, hi := uint32(hkey>>32), uint32(hkey)
	for i := uint64(0); i < vec.d; i++ {
		if ce := uint64(atomic.LoadUint32(&vec.buf[vecpos(lo, hi, vec.w, i)])); r == 0 || r > ce {
			r = ce
		}
	}
	return
}

func (vec *cnvector32) decay(ctx context.Context, factor float64) error {
	for i := 0; i < len(vec.buf); i++ {
		var ok bool
		var j uint64
		for j = 0; j < vec.lim+1; j++ {
			o := atomic.LoadUint32(&vec.buf[i])
			n := uint32(float64(o) * factor)
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				if ok = atomic.CompareAndSwapUint32(&vec.buf[i], o, n); ok {
					break
				}
			}
		}
		if !ok {
			return pbtk.ErrWriteLimitExceed
		}
	}
	return nil
}

func (vec *cnvector32) reset() {
	for i := uint64(0); i < uint64(len(vec.buf)); i++ {
		atomic.StoreUint32(&vec.buf[i], 0)
	}
}

func (vec *cnvector32) readFrom(r io.Reader) (n int64, err error) {
	var (
		buf [24]byte
		m   int
	)
	m, err = r.Read(buf[:])
	n += int64(m)
	if err != nil {
		return
	}

	if binary.LittleEndian.Uint64(buf[0:8]) != dumpSignatureConcurrent32 {
		err = pbtk.ErrInvalidSignature
		return
	}
	if binary.LittleEndian.Uint64(buf[8:16]) != dumpVersionConcurrent32 {
		err = pbtk.ErrVersionMismatch
		return
	}
	vec.flags = bitset.Bitset64(binary.LittleEndian.Uint64(buf[16:24]))

	for i := 0; i < len(vec.buf); i++ {
		m, err = r.Read(buf[:4])
		n += int64(m)
		if err != nil {
			return
		}
		v := binary.LittleEndian.Uint32(buf[:4])
		atomic.StoreUint32(&vec.buf[i], v)
	}
	return
}

func (vec *cnvector32) writeTo(w io.Writer) (n int64, err error) {
	var (
		buf [24]byte
		m   int
	)
	binary.LittleEndian.PutUint64(buf[0:8], dumpSignatureConcurrent32)
	binary.LittleEndian.PutUint64(buf[8:16], dumpVersionConcurrent32)
	binary.LittleEndian.PutUint64(buf[16:24], uint64(vec.flags))
	m, err = w.Write(buf[:])
	n += int64(m)
	if err != nil {
		return
	}

	for i := 0; i < len(vec.buf); i++ {
		v := atomic.LoadUint32(&vec.buf[i])
		binary.LittleEndian.PutUint32(buf[0:4], v)
		m, err = w.Write(buf[:4])
		n += int64(m)
		if err != nil {
			return
		}
	}
	return
}

func newConcurrentVector32(d, w, lim uint64, flags bitset.Bitset64) *cnvector32 {
	return &cnvector32{
		basevec: basevec{d: d, w: w, flags: flags},
		lim:     lim,
		buf:     make([]uint32, d*w),
	}
}
