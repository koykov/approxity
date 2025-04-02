package tinylfu

import (
	"context"
	"math"
	"sync"
	"sync/atomic"
	"time"

	"github.com/koykov/pbtk"
	"github.com/koykov/pbtk/frequency"
	"github.com/koykov/pbtk/frequency/cmsketch"
)

const flagLFU = 1

type estimator[T pbtk.Hashable] struct {
	conf   *Config
	est    frequency.Estimator[T]
	dec    frequency.Decayer
	once   sync.Once
	cancel context.CancelFunc // main stop func
	timer  *time.Timer        // timer reached notifier
	c      uint64             // counter of added items
	cntr   chan struct{}      // counter reached notifier
	svc    uint32             // decay running flag
	lt     int64              // last decay timestamp

	err error
}

func NewEstimator[T pbtk.Hashable](conf *Config) (frequency.Estimator[T], error) {
	if conf == nil {
		return nil, pbtk.ErrInvalidConfig
	}
	conf.WithFlag(flagLFU, true)
	cms, err := cmsketch.NewEstimator[T](&conf.Config)
	if err != nil {
		return nil, err
	}
	dec := any(cms).(frequency.Decayer)
	e := &estimator[T]{
		conf: conf.copy(),
		est:  cms,
		dec:  dec,
	}
	if e.once.Do(e.init); e.err != nil {
		return nil, e.err
	}
	return e, nil
}

func (e *estimator[T]) Add(key T) error                   { return e.checkCntr(e.est.Add(key)) }
func (e *estimator[T]) AddN(key T, n uint64) error        { return e.checkCntr(e.est.AddN(key, n)) }
func (e *estimator[T]) HAdd(hkey uint64) error            { return e.checkCntr(e.est.HAdd(hkey)) }
func (e *estimator[T]) HAddN(hkey uint64, n uint64) error { return e.checkCntr(e.est.HAddN(hkey, n)) }

func (e *estimator[T]) checkCntr(err error) error {
	if err != nil {
		return err
	}
	if atomic.AddUint64(&e.c, 1) == e.conf.DecayLimit {
		e.cntr <- struct{}{}
	}
	return nil
}

func (e *estimator[T]) Close() error {
	e.cancel()
	return nil
}

func (e *estimator[T]) init() {
	if e.conf.DecayLimit == 0 {
		e.conf.DecayLimit = defaultDecayLimit
	}
	if e.conf.DecayInterval == 0 {
		e.conf.DecayInterval = defaultDecayInterval
	}
	if e.conf.DecayFactor == 0 {
		e.conf.DecayFactor = defaultDecayFactor
	}
	if e.conf.SoftDecayFactor == 0 {
		e.conf.SoftDecayFactor = defaultSoftDecayFactor
	}
	if e.conf.ForceDecayNotifier == nil {
		e.conf.ForceDecayNotifier = dummyForceDecayNotifier{}
	}
	if e.conf.Concurrent == nil {
		// only concurrent CMS allowed due to async decay
		e.conf.Concurrent = &cmsketch.ConcurrentConfig{}
	}

	// counter
	e.c = math.MaxUint64
	e.cntr = make(chan struct{})
	// timer
	e.timer = time.NewTimer(e.conf.DecayInterval)

	ctx, cancel := context.WithCancel(context.Background())
	e.cancel = cancel
	go e.watch(ctx)
}

func (e *estimator[T]) watch(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			e.err = pbtk.ErrClosed
			e.timer.Stop()
			close(e.cntr)
			return
		case <-e.conf.ForceDecayNotifier.Notify():
			e.decay(ctx)
		case <-e.timer.C:
			e.decay(ctx)
		case <-e.cntr:
			e.decay(ctx)
		}
	}
}

func (e *estimator[T]) decay(ctx context.Context) {
	if !atomic.CompareAndSwapUint32(&e.svc, 0, 1) {
		return
	}
	defer atomic.StoreUint32(&e.svc, 0)

	factor := e.conf.DecayFactor
	{
		// try soft decay
		var interval, counter bool
		if lt := atomic.LoadInt64(&e.lt); lt > 0 {
			left := time.Now().Sub(time.Unix(0, lt))
			interval = left > 0 && left < e.conf.DecayInterval/2
		}
		c := atomic.LoadUint64(&e.c)
		counter = c > 0 && c < e.conf.DecayLimit/2
		if interval || counter {
			factor = e.conf.SoftDecayFactor
		}
	}

	e.timer.Reset(e.conf.DecayInterval)
	atomic.StoreUint64(&e.c, 0)
	atomic.StoreInt64(&e.lt, time.Now().UnixNano())
	_ = e.dec.Decay(ctx, factor)
}
