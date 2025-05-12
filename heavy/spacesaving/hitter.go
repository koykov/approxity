package spacesaving

import (
	"sync"

	"github.com/koykov/pbtk"
	"github.com/koykov/pbtk/heavy"
)

type hitter[T pbtk.Hashable] struct {
	conf *Config
	once sync.Once

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

func (h *hitter[T]) Add(k T) error {
	// todo implement me
	return nil
}

func (h *hitter[T]) Hits() []heavy.Hit[T] {
	// todo implement me
	return nil
}

func (h *hitter[T]) AppendHits(dst []heavy.Hit[T]) []heavy.Hit[T] {
	// todo implement me
	return dst
}

func (h *hitter[T]) Reset() {
	// todo implement me
}

func (h *hitter[T]) init() {
	// todo implement me
}
