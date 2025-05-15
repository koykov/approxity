package lossy

import (
	"math"
	"slices"
	"sync"

	"github.com/koykov/pbtk"
	"github.com/koykov/pbtk/heavy"
)

type hitter[T pbtk.Hashable] struct {
	pbtk.Base[T]
	conf    *Config
	buckets []*bucket[T]
	w       uint64
	once    sync.Once

	err error
}

func NewHitter[T pbtk.Hashable](conf *Config) (heavy.Hitter[T], error) {
	if conf == nil {
		return nil, pbtk.ErrInvalidConfig
	}

	h := &hitter[T]{conf: conf}
	if h.once.Do(h.init); h.err != nil {
		return nil, h.err
	}

	return h, nil
}

func (h *hitter[T]) Add(key T) error {
	if h.once.Do(h.init); h.err != nil {
		return h.mw().Add(h.err)
	}
	hkey, err := h.Hash(h.conf.Hasher, key)
	if err != nil {
		return h.mw().Add(err)
	}
	bi := hkey % h.conf.Buckets
	h.buckets[bi].add(key, hkey)
	return h.mw().Add(nil)
}

func (h *hitter[T]) Hits() []heavy.Hit[T] {
	if h.once.Do(h.init); h.err != nil {
		return nil
	}
	buf := make([]heavy.Hit[T], 0, h.conf.Buckets*h.w)
	return h.appendHits(buf)
}

func (h *hitter[T]) AppendHits(dst []heavy.Hit[T]) []heavy.Hit[T] {
	if h.once.Do(h.init); h.err != nil {
		return dst
	}
	return h.appendHits(dst)
}

func (h *hitter[T]) appendHits(dst []heavy.Hit[T]) []heavy.Hit[T] {
	for i := 0; i < len(h.buckets); i++ {
		dst = h.buckets[i].appendHits(dst)
	}
	if len(dst) == 0 {
		return dst
	}
	slices.SortFunc(dst, func(a, b heavy.Hit[T]) int {
		// reverse order
		switch {
		case a.Rate > b.Rate:
			return -1
		case a.Rate < b.Rate:
			return 1
		}
		return 0
	})
	h.mw().Hits(dst[0].Rate, dst[len(dst)-1].Rate)
	return dst
}

func (h *hitter[T]) Reset() {
	if h.once.Do(h.init); h.err != nil {
		return
	}
	for i := 0; i < len(h.buckets); i++ {
		h.buckets[i].reset()
	}
}

func (h *hitter[T]) init() {
	if h.conf.Hasher == nil {
		h.err = pbtk.ErrNoHasher
		return
	}
	if h.conf.Epsilon == 0 {
		h.err = ErrZeroEpsilon
		return
	}
	if h.conf.Epsilon >= h.conf.Support {
		h.err = ErrBadEpsilon
		return
	}
	if h.conf.Buckets == 0 {
		h.conf.Buckets = defaultBuckets
	}
	if h.conf.MetricsWriter == nil {
		h.conf.MetricsWriter = &heavy.DummyMetricsWriter{}
	}
	h.w = uint64(math.Ceil(1 / h.conf.Epsilon))
	h.buckets = make([]*bucket[T], h.conf.Buckets)
	for i := range h.buckets {
		h.buckets[i] = &bucket[T]{
			e:    h.conf.Epsilon,
			s:    h.conf.Support,
			w:    h.w,
			keys: make(map[uint64]*tuple[T], h.w),
		}
	}
}

func (h *hitter[T]) mw() heavy.MetricsWriter {
	return h.conf.MetricsWriter
}
