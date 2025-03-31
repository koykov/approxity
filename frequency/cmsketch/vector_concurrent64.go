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
	dumpSignatureConcurrent64 = 0xABF100F41194C630
	dumpVersionConcurrent64   = 1.0
)

// 64-bit version of concurrent vector implementation.
type cnvector64 struct {
	basevec
	lim uint64
	buf []uint64
}

func (vec *cnvector64) add(hkey, delta uint64) error {
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

func (vec *cnvector64) addClassic(lo, hi uint32, delta uint64) error {
	for i := uint64(0); i < vec.d; i++ {
		pos := i*vec.w + uint64(lo+hi*uint32(i))%vec.w
		var j uint64
		for j = 0; j < vec.lim+1; j++ {
			o := atomic.LoadUint64(&vec.buf[pos])
			n := o + delta
			if atomic.CompareAndSwapUint64(&vec.buf[pos], o, n) {
				break
			}
		}
		if j == vec.lim+1 {
			return pbtk.ErrWriteLimitExceed
		}
	}
	return nil
}

func (vec *cnvector64) addCU(lo, hi uint32, delta uint64) error {
	var mn uint64 = math.MaxUint64
	for i := uint64(0); i < vec.d; i++ {
		mn = min(mn, atomic.LoadUint64(&vec.buf[vecpos(lo, hi, vec.w, i)]))
	}
	for i := uint64(0); i < vec.d; i++ {
		pos := vecpos(lo, hi, vec.w, i)
		if atomic.LoadUint64(&vec.buf[pos]) == mn {
			var j uint64
			for j = 0; j < vec.lim+1; j++ {
				o := atomic.LoadUint64(&vec.buf[pos])
				n := o + delta
				if atomic.CompareAndSwapUint64(&vec.buf[pos], o, n) {
					break
				}
			}
			if j == vec.lim+1 {
				return pbtk.ErrWriteLimitExceed
			}
		}
	}
	return nil
}

func (vec *cnvector64) addLFU(lo, hi uint32, delta uint64) error {
	return vec.addClassic(lo, hi, delta)
}

func (vec *cnvector64) estimate(hkey uint64) (r uint64) {
	lo, hi := uint32(hkey>>32), uint32(hkey)
	for i := uint64(0); i < vec.d; i++ {
		if ce := atomic.LoadUint64(&vec.buf[vecpos(lo, hi, vec.w, i)]); r == 0 || r > ce {
			r = ce
		}
	}
	return
}

func (vec *cnvector64) decay(ctx context.Context, factor float64) error {
	for i := 0; i < len(vec.buf); i++ {
		o := atomic.LoadUint64(&vec.buf[i])
		n := uint64(float64(o) * factor)
		var j uint64
		for j = 0; j < vec.lim+1; j++ {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				if atomic.CompareAndSwapUint64(&vec.buf[i], o, n) {
					return nil
				}
			}
		}
		if j == vec.lim+1 {
			return pbtk.ErrWriteLimitExceed
		}
	}
	return nil
}

func (vec *cnvector64) reset() {
	for i := uint64(0); i < uint64(len(vec.buf)); i++ {
		atomic.StoreUint64(&vec.buf[i], 0)
	}
}

func (vec *cnvector64) readFrom(r io.Reader) (n int64, err error) {
	var (
		buf [24]byte
		m   int
	)
	m, err = r.Read(buf[:])
	n += int64(m)
	if err != nil {
		return
	}

	if binary.LittleEndian.Uint64(buf[0:8]) != dumpSignatureConcurrent64 {
		err = pbtk.ErrInvalidSignature
		return
	}
	if binary.LittleEndian.Uint64(buf[8:16]) != dumpVersionConcurrent64 {
		err = pbtk.ErrVersionMismatch
		return
	}
	vec.flags = bitset.Bitset64(binary.LittleEndian.Uint64(buf[16:24]))

	for i := 0; i < len(vec.buf); i++ {
		m, err = r.Read(buf[:8])
		n += int64(m)
		if err != nil {
			return
		}
		v := binary.LittleEndian.Uint64(buf[:8])
		atomic.StoreUint64(&vec.buf[i], v)
	}
	return
}

func (vec *cnvector64) writeTo(w io.Writer) (n int64, err error) {
	var (
		buf [24]byte
		m   int
	)
	binary.LittleEndian.PutUint64(buf[0:8], dumpSignatureConcurrent64)
	binary.LittleEndian.PutUint64(buf[8:16], dumpVersionConcurrent64)
	binary.LittleEndian.PutUint64(buf[16:24], uint64(vec.flags))
	m, err = w.Write(buf[:])
	n += int64(m)
	if err != nil {
		return
	}

	for i := 0; i < len(vec.buf); i++ {
		v := atomic.LoadUint64(&vec.buf[i])
		binary.LittleEndian.PutUint64(buf[0:8], v)
		m, err = w.Write(buf[:8])
		n += int64(m)
		if err != nil {
			return
		}
	}
	return
}

func newConcurrentVector64(d, w, lim uint64, flags bitset.Bitset64) *cnvector64 {
	return &cnvector64{
		basevec: basevec{d: d, w: w, flags: flags},
		lim:     lim,
		buf:     make([]uint64, d*w),
	}
}
