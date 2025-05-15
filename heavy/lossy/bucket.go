package lossy

import (
	"sync"

	"github.com/koykov/pbtk"
	"github.com/koykov/pbtk/heavy"
)

type bucket[T pbtk.Hashable] struct {
	e, s float64
	w    uint64
	mux  sync.RWMutex
	n    uint64
	keys map[uint64]*tuple[T]
}

type tuple[T pbtk.Hashable] struct {
	key   T
	f     float64
	delta float64
}

func (b *bucket[T]) add(key T, hkey uint64) {
	b.n++
	idx := b.n / b.w
	b.mux.Lock()
	defer b.mux.Unlock()
	t, ok := b.keys[hkey]
	if ok {
		t.f++
	} else {
		t = &tuple[T]{
			key:   key,
			f:     1,
			delta: float64(idx - 1),
		}
		b.keys[hkey] = t
	}
	if b.n%b.w == 0 {
		for k, t := range b.keys {
			if t.f+t.delta <= float64(idx) {
				delete(b.keys, k)
			}
		}
	}
}

func (b *bucket[T]) appendHits(dst []heavy.Hit[T]) []heavy.Hit[T] {
	b.mux.RLock()
	defer b.mux.RUnlock()
	for _, t := range b.keys {
		if t.f >= 1-b.e*float64(b.n) {
			dst = append(dst, heavy.Hit[T]{
				Key:  t.key,
				Rate: t.f/float64(b.n) + b.s,
			})
		}
	}
	return dst
}

func (b *bucket[T]) reset() {
	b.mux.Lock()
	defer b.mux.Unlock()
	b.n = 0
	for k := range b.keys {
		delete(b.keys, k)
	}
}
