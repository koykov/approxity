package tinylfu

import (
	"encoding/binary"
	"io"
	"math"
	"sync/atomic"

	"github.com/koykov/pbtk"
)

const (
	cnvecDumpSignature = 0x813b6cd70a883800
	cnvecDumpVersion   = 1.0
)

type cnvec struct {
	*basevec
	lim uint64
}

func (vec *cnvec) set(pos, n uint64, dtime uint32) error {
	for i := uint64(0); i < vec.lim+1; i++ {
		val := atomic.LoadUint64(&vec.buf[pos])
		newVal := vec.recalc(val, n, dtime)
		if atomic.CompareAndSwapUint64(&vec.buf[pos], val, newVal) {
			return nil
		}
	}
	return pbtk.ErrWriteLimitExceed
}

func (vec *cnvec) get(pos uint64, stime, now uint32) uint32 {
	val := atomic.LoadUint64(&vec.buf[pos])
	return vec.estimate(val, stime, now)
}

func (vec *cnvec) reset() {
	for i := 0; i < len(vec.buf); i++ {
		atomic.StoreUint64(&vec.buf[i], 0)
	}
}

func (vec *cnvec) readFrom(r io.Reader) (n int64, err error) {
	var (
		buf [64]byte
		m   int
	)
	m, err = r.Read(buf[:])
	n += int64(m)
	if err != nil {
		return n, err
	}

	sign, ver, dtimeMin, tau, decayMin, bufsz, exptabsz, lim := binary.LittleEndian.Uint64(buf[0:8]),
		binary.LittleEndian.Uint64(buf[8:16]), binary.LittleEndian.Uint64(buf[16:24]),
		binary.LittleEndian.Uint64(buf[24:32]), binary.LittleEndian.Uint64(buf[32:40]),
		binary.LittleEndian.Uint64(buf[40:48]), binary.LittleEndian.Uint64(buf[48:56]),
		binary.LittleEndian.Uint64(buf[56:64])

	if sign != syncvecDumpSignature {
		return n, pbtk.ErrInvalidSignature
	}
	if ver != math.Float64bits(syncvecDumpVersion) {
		return n, pbtk.ErrVersionMismatch
	}
	vec.dtimeMin, vec.tau = dtimeMin, tau
	vec.decayMin = math.Float64frombits(decayMin)
	vec.lim = lim

	if uint64(len(vec.buf)) < bufsz {
		vec.buf = make([]uint64, bufsz)
	}
	vec.buf = vec.buf[:bufsz]
	for i := uint64(0); i < bufsz; i++ {
		m, err = r.Read(buf[0:8])
		n += int64(m)
		if err != nil {
			return n, err
		}
		atomic.StoreUint64(&vec.buf[i], binary.LittleEndian.Uint64(buf[0:8]))
	}

	if uint64(len(vec.exptab)) < exptabsz {
		vec.exptab = make([]float64, exptabsz)
	}
	vec.exptab = vec.exptab[:exptabsz]
	for i := uint64(0); i < exptabsz; i++ {
		m, err = r.Read(buf[0:8])
		n += int64(m)
		if err == io.EOF {
			err = nil
		}
		if err != nil {
			return n, err
		}
		vec.exptab[i] = math.Float64frombits(binary.LittleEndian.Uint64(buf[0:8]))
	}
	if err == io.EOF {
		err = nil
	}

	return
}

func (vec *cnvec) writeTo(w io.Writer) (n int64, err error) {
	var (
		buf [64]byte
		m   int
	)
	binary.LittleEndian.PutUint64(buf[0:8], syncvecDumpSignature)
	binary.LittleEndian.PutUint64(buf[8:16], math.Float64bits(syncvecDumpVersion))
	binary.LittleEndian.PutUint64(buf[16:24], vec.dtimeMin)
	binary.LittleEndian.PutUint64(buf[24:32], vec.tau)
	binary.LittleEndian.PutUint64(buf[32:40], math.Float64bits(vec.decayMin))
	binary.LittleEndian.PutUint64(buf[40:48], uint64(len(vec.buf)))
	binary.LittleEndian.PutUint64(buf[48:56], vec.exptabsz)
	binary.LittleEndian.PutUint64(buf[56:64], vec.lim)
	m, err = w.Write(buf[:])
	n += int64(m)
	if err != nil {
		return int64(m), err
	}

	for i := 0; i < len(vec.buf); i++ {
		binary.LittleEndian.PutUint64(buf[0:8], atomic.LoadUint64(&vec.buf[i]))
		m, err = w.Write(buf[:8])
		n += int64(m)
		if err != nil {
			return int64(m), err
		}
	}
	for i := 0; i < len(vec.exptab); i++ {
		binary.LittleEndian.PutUint64(buf[0:8], math.Float64bits(vec.exptab[i]))
		m, err = w.Write(buf[:8])
		n += int64(m)
		if err != nil {
			return int64(m), err
		}
	}
	return
}

func newConcurrentVector(sz, lim uint64, ewma *EWMA) vector {
	vec := &cnvec{
		basevec: &basevec{
			buf:      make([]uint64, sz),
			dtimeMin: ewma.MinDeltaTime,
			tau:      ewma.Tau,
			exptabsz: ewma.ExpTableSize,
		},
		lim: lim,
	}
	vec.basevec.init()
	return vec
}
