package tinylfu

import (
	"io"
	"math"
	"sync"

	"github.com/koykov/pbtk"
	"github.com/koykov/pbtk/frequency"
)

type estimator[T pbtk.Hashable] struct {
	pbtk.Base[T]
	conf  *Config
	once  sync.Once
	stime uint32 // start time
	w, d  uint64
	vec   vector

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
	now := e.now()
	timeDeltaNew := now - e.stime
	for i := uint64(0); i < e.d; i++ {
		pos := i*e.w + hkey%e.w
		if err := e.vec.set(pos, n, timeDeltaNew); err != nil {
			return err
		}
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
		val := e.vec.get(pos, e.stime, now)
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
	e.vec.reset()
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
	if e.conf.EWMA.Tau == 0 {
		e.conf.EWMA.Tau = defaultTau
	}
	if e.conf.EWMA.MinDeltaTime == 0 {
		e.conf.EWMA.MinDeltaTime = defaultMinDeltaTime
	}
	if e.conf.EWMA.TimePrecision == 0 {
		e.conf.EWMA.TimePrecision = defaultTimePrecision
	}
	if e.conf.EWMA.ExpTableSize == 0 {
		e.conf.EWMA.ExpTableSize = defaultExpTableSize
	}
	if e.conf.Clock == nil {
		e.conf.Clock = nativeClock{}
	}
	if e.conf.MetricsWriter == nil {
		e.conf.MetricsWriter = frequency.DummyMetricsWriter{}
	}

	e.w, e.d = optimalWD(e.conf.Confidence, e.conf.Epsilon)
	if e.conf.Concurrent != nil {
		e.vec = newConcurrentVector(e.w*e.d, e.conf.Concurrent.WriteAttemptsLimit, &e.conf.EWMA)
	} else {
		e.vec = newVector(e.w*e.d, &e.conf.EWMA)
	}
	e.stime = e.now()
}

func (e *estimator[T]) now() uint32 {
	return uint32(e.conf.Clock.Now().UnixNano() / int64(e.conf.EWMA.TimePrecision))
}
