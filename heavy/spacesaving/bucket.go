package spacesaving

import (
	"sync"

	"github.com/koykov/pbtk"
	"github.com/koykov/pbtk/heavy"
)

type bucket[T pbtk.Hashable] struct {
	k    uint64
	a    float64
	mux  sync.RWMutex
	keys map[uint64]uint64
	buf  []tuple[T]
}

type tuple[T pbtk.Hashable] struct {
	key  T
	hkey uint64
	rate float64
}

func (b *bucket[T]) add(key T, hkey, n uint64) {
	b.mux.Lock()
	defer b.mux.Unlock()
	if i, ok := b.keys[hkey]; ok {
		b.buf[i].rate += b.tryEWMA(b.buf[i].rate, n)
		return
	}
	if uint64(len(b.buf)) < b.k {
		t := tuple[T]{
			key:  key,
			rate: 1,
			hkey: hkey,
		}
		b.buf = append(b.buf, t)
		b.keys[hkey] = uint64(len(b.buf) - 1)
		return
	}
	_ = b.buf[len(b.buf)-1]
	mt, mi := &b.buf[0], 0
	for i := 1; i < len(b.buf); i++ {
		if b.buf[i].rate < mt.rate {
			mt, mi = &b.buf[i], i
		}
	}
	delete(b.keys, mt.hkey)
	mt.key = key
	mt.rate = b.tryEWMA(mt.rate, n)
	mt.hkey = hkey
	b.keys[hkey] = uint64(mi)
}

func (b *bucket[T]) tryEWMA(val float64, n uint64) float64 {
	if b.a == 0 {
		return val + float64(n)
	}
	return b.a*float64(n) + (1-b.a)*val
}

func (b *bucket[T]) appendHits(dst []heavy.Hit[T]) []heavy.Hit[T] {
	b.mux.RLock()
	defer b.mux.RUnlock()
	if len(b.buf) == 0 {
		return dst
	}
	_ = b.buf[len(b.buf)-1]
	for i := 0; i < len(b.buf); i++ {
		t := &b.buf[i]
		dst = append(dst, heavy.Hit[T]{
			Key:  t.key,
			Rate: t.rate,
		})
	}
	return dst
}

func (b *bucket[T]) reset() {
	b.mux.Lock()
	defer b.mux.Unlock()
	for k := range b.keys {
		delete(b.keys, k)
	}
	b.buf = b.buf[:0]
}
