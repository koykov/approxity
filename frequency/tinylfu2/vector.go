package tinylfu

import (
	"io"
	"math"
)

type vector interface {
	set(pos, n uint64, dtime uint32) error
	get(pos uint64, stime, now uint32) float64
	reset()
	readFrom(r io.Reader) (int64, error)
	writeTo(w io.Writer) (int64, error)
}

type basevec struct {
	buf    []uint64
	exptab []float64

	dtimeMin, tau uint64
	decayMin      float64
	exptabsz      uint64
}

func (vec *basevec) encode(tdelta uint32, n uint32) uint64 {
	return (uint64(tdelta) << 32) | uint64(n)
}

func (vec *basevec) decode(val uint64) (tdelta uint32, n uint32) {
	tdelta = uint32(val >> 32)
	n = uint32(val & 0xFFFFFFFF)
	return
}

func (vec *basevec) recalc(val, n uint64, dtimeNew uint32) uint64 {
	dtimeOld, valOld := vec.decode(val)
	var valNew float64
	if valOld == 0 && dtimeOld == 0 {
		// first addition
		valNew = float64(n)
	} else {
		dtime := dtimeNew - dtimeOld
		var decay float64
		if uint64(dtime) < vec.dtimeMin {
			// special case - update item before minDeltaTime
			decay = vec.decayMin
			valNew = float64(valOld) + float64(n)*(1-decay)
		} else {
			// regular case - update item after minDeltaTime since addition
			decay = vec.exp(dtime) // e^(-Δt/τ)
			valNew = float64(valOld)*decay + float64(n)*(1-decay)
		}
	}
	return vec.encode(dtimeNew, uint32(valNew))
}

func (vec *basevec) estimate(val uint64, stime, now uint32) float64 {
	timeDeltaOld, valOld := vec.decode(val)
	if valOld == 0 && timeDeltaOld == 0 {
		return math.MaxUint32
	}
	timeDelta := now - stime - timeDeltaOld
	decay := vec.exp(timeDelta) // e^(-Δt/τ)
	return float64(valOld) * decay
}

func (vec *basevec) exp(dtime uint32) float64 {
	if uint64(dtime) >= vec.exptabsz {
		return math.Exp(-float64(dtime) / float64(vec.tau))
	}
	return vec.exptab[dtime]
}

func (vec *basevec) init() {
	vec.decayMin = math.Exp(-float64(vec.dtimeMin) / float64(vec.tau))
	vec.exptab = make([]float64, vec.exptabsz)
	for i := uint64(0); i < vec.exptabsz; i++ {
		vec.exptab[i] = math.Exp(-float64(i) / float64(vec.tau))
	}
}
