package spacesaving

import (
	"slices"
	"sync"

	"github.com/koykov/pbtk"
	"github.com/koykov/pbtk/heavy"
)

type hitter[T pbtk.Hashable] struct {
	pbtk.Base[T]
	conf    *Config
	once    sync.Once
	buckets []*bucket[T]

	err error
}

func NewHitter[T pbtk.Hashable](conf *Config) (heavy.Hitter[T], error) {
	if conf == nil {
		return nil, pbtk.ErrInvalidConfig
	}
	h := &hitter[T]{conf: conf.copy()}
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
	h.buckets[bi].add(key, hkey, 1)
	return h.mw().Add(nil)
}

func (h *hitter[T]) Hits() []heavy.Hit[T] {
	if h.once.Do(h.init); h.err != nil {
		return nil
	}
	buf := make([]heavy.Hit[T], 0, h.conf.K*h.conf.Buckets)
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
	return dst[:h.conf.K]
}

func (h *hitter[T]) Reset() {
	if h.once.Do(h.init); h.err != nil {
		return
	}
	h.mw().Reset()
	for i := 0; i < len(h.buckets); i++ {
		h.buckets[i].reset()
	}
}

func (h *hitter[T]) init() {
	if h.conf.Hasher == nil {
		h.err = pbtk.ErrNoHasher
		return
	}
	if h.conf.K == 0 {
		h.err = heavy.ErrZeroK
		return
	}
	if h.conf.Buckets == 0 {
		h.conf.Buckets = defaultBuckets
	}
	if h.conf.MetricsWriter == nil {
		h.conf.MetricsWriter = &heavy.DummyMetricsWriter{}
	}
	for i := uint64(0); i < h.conf.Buckets; i++ {
		b := &bucket[T]{
			k:    h.conf.K,
			a:    h.conf.EWMA.Alpha,
			keys: make(map[uint64]uint64, h.conf.K),
			buf:  make([]tuple[T], 0, h.conf.K),
		}
		h.buckets = append(h.buckets, b)
	}
}

func (h *hitter[T]) mw() heavy.MetricsWriter {
	return h.conf.MetricsWriter
}
