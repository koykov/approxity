package misragries

import (
	"sync"

	"github.com/koykov/pbtk"
	"github.com/koykov/pbtk/heavy"
)

type bucket[T pbtk.Hashable] struct {
	k    uint64
	mux  sync.RWMutex
	keys map[uint64]uint64
	buf  []tuple[T]
}

type tuple[T pbtk.Hashable] struct {
	key  T
	hkey uint64
	rate float64
}

func (b *bucket[T]) add(key T, hkey uint64) {
	b.mux.Lock()
	defer b.mux.Unlock()
	if i, ok := b.keys[hkey]; ok {
		b.buf[i].rate++
		return
	}
	if uint64(len(b.buf)) < b.k {
		b.buf = append(b.buf, tuple[T]{
			key:  key,
			hkey: hkey,
			rate: 1,
		})
		b.keys[hkey] = uint64(len(b.buf) - 1)
		return
	}
	_ = b.buf[len(b.buf)-1]
	pos := -1
	for i := uint64(0); i < b.k; i++ {
		if b.buf[i].rate > 0 {
			b.buf[i].rate--
		}
		if b.buf[i].rate == 0 && pos == -1 {
			pos = int(i)
		}
	}
	if pos != -1 {
		delete(b.keys, b.buf[pos].hkey)
		b.buf[pos].key = key
		b.buf[pos].hkey = hkey
		b.buf[pos].rate = 1
		b.keys[hkey] = uint64(pos)
	}
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
		if t.rate == 0 {
			continue
		}
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
