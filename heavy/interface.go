package heavy

import (
	"github.com/koykov/pbtk"
)

type Hitter[T pbtk.Hashable] interface {
	Add(key T) error
	Hits() []Hit[T]
	AppendHits(dst []Hit[T]) []Hit[T]
	Reset()
}

type Hit[T pbtk.Hashable] struct {
	Key  T
	Rate float64
}

func (h *Hit[T]) Freq() float64 {
	return h.Rate
}
