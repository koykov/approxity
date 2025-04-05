package tinylfu

import (
	"io"
	"math"
	"sync"
	"unsafe"

	"github.com/koykov/openrt"
	"github.com/koykov/pbtk"
	"github.com/koykov/pbtk/frequency"
)

type estimator[T pbtk.Hashable] struct {
	pbtk.Base[T]
	conf  *Config
	once  sync.Once
	stime uint32 // start time
	w, d  uint64
	vec   []uint64

	err error
}

func NewEstimator[T pbtk.Hashable](conf *Config) (frequency.Estimator[T], error) {
	if conf == nil {
		return nil, pbtk.ErrInvalidConfig
	}
	e := &estimator[T]{conf: conf.copy()}
	if e.once.Do(e.init); e.err != nil {
		return nil, e.err
	}
	return e, nil
}

func (e *estimator[T]) Add(key T) error {
	return e.AddN(key, 1)
}

func (e *estimator[T]) AddN(key T, n uint64) error {
	if e.once.Do(e.init); e.err != nil {
		return e.err
	}
	hkey, err := e.Hash(e.conf.Hasher, key)
	if err != nil {
		return err
	}
	return e.HAddN(hkey, n)
}

func (e *estimator[T]) HAdd(hkey uint64) error {
	return e.HAddN(hkey, 1)
}

func (e *estimator[T]) HAddN(hkey uint64, n uint64) error {
	if e.once.Do(e.init); e.err != nil {
		return e.err
	}
	timeDeltaNew := e.now() - e.stime
	for i := uint64(0); i < e.d; i++ {
		pos := i*e.w + hkey%e.w
		timeDeltaOld, valOld := decode(e.vec[pos])
		var valNew uint32
		if valOld == 0 && timeDeltaOld == 0 {
			valNew = uint32(n)
		} else {
			timeDelta := timeDeltaNew - timeDeltaOld
			decay := math.Exp(-float64(timeDelta) / float64(e.conf.EWMA.Tau)) // e^(-Δt/τ)
			valNew = uint32(float64(valOld)*decay + float64(n)*(1-decay))
		}
		e.vec[pos] = encode(timeDeltaNew, valNew)
	}
	return nil
}

func (e *estimator[T]) Estimate(key T) uint64 {
	if e.once.Do(e.init); e.err != nil {
		return 0
	}
	hkey, err := e.Hash(e.conf.Hasher, key)
	if err != nil {
		return 0
	}
	return e.HEstimate(hkey)
}

func (e *estimator[T]) HEstimate(hkey uint64) uint64 {
	if e.once.Do(e.init); e.err != nil {
		return 0
	}
	now := e.now()
	minVal := uint32(math.MaxUint32)
	for i := uint64(0); i < e.d; i++ {
		pos := i*e.w + hkey%e.w
		timeDeltaOld, valOld := decode(e.vec[pos])
		if valOld == 0 && timeDeltaOld == 0 {
			continue
		}
		timeDelta := now - e.stime - timeDeltaOld
		decay := math.Exp(-float64(timeDelta) / float64(e.conf.EWMA.Tau)) // e^(-Δt/τ)
		val := uint32(float64(valOld) * decay)
		if val < minVal {
			minVal = val
		}
	}
	if minVal == math.MaxUint32 {
		minVal = 0
	}
	return uint64(minVal)
}

func (e *estimator[T]) Reset() {
	if e.once.Do(e.init); e.err != nil {
		return
	}
	openrt.MemclrUnsafe(unsafe.Pointer(&e.vec[0]), len(e.vec)*8)
}

func (e *estimator[T]) ReadFrom(r io.Reader) (int64, error) {
	if e.once.Do(e.init); e.err != nil {
		return 0, e.err
	}
	// todo implement me
	return 0, nil
}

func (e *estimator[T]) WriteTo(w io.Writer) (int64, error) {
	if e.once.Do(e.init); e.err != nil {
		return 0, e.err
	}
	// todo implement me
	return 0, nil
}

func (e *estimator[T]) init() {
	if e.conf.Hasher == nil {
		e.err = pbtk.ErrNoHasher
		return
	}
	if e.conf.Confidence <= 0 || e.conf.Confidence >= 1 {
		e.err = frequency.ErrInvalidConfidence
		return
	}
	if e.conf.Epsilon <= 0 || e.conf.Epsilon >= 1 {
		e.err = frequency.ErrInvalidEpsilon
		return
	}
	if e.conf.Clock == nil {
		e.conf.Clock = nativeClock{}
	}
	if e.conf.MetricsWriter == nil {
		e.conf.MetricsWriter = frequency.DummyMetricsWriter{}
	}

	e.w, e.d = optimalWD(e.conf.Confidence, e.conf.Epsilon)
	e.vec = make([]uint64, e.w*e.d)
	e.stime = e.now()
}

func (e *estimator[T]) now() uint32 {
	return e.conf.Clock.UNow32()
}
