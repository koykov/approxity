package tinylfu

import "time"

type timer interface {
	C() <-chan time.Time
	Stop() bool
	Reset(duration time.Duration) bool
}

type nativeTimer struct {
	t *time.Timer
}

func newNativeTimer(d time.Duration) timer {
	return &nativeTimer{t: time.NewTimer(d)}
}

func (t *nativeTimer) C() <-chan time.Time        { return t.t.C }
func (t *nativeTimer) Stop() bool                 { return t.t.Stop() }
func (t *nativeTimer) Reset(d time.Duration) bool { return t.t.Reset(d) }

type stuckTimer struct{}

func (t *stuckTimer) C() <-chan time.Time        { return nil }
func (t *stuckTimer) Stop() bool                 { return true }
func (t *stuckTimer) Reset(_ time.Duration) bool { return true }
