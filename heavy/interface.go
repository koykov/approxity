package heavy

import (
	"github.com/koykov/pbtk"
)

type Hit[T pbtk.Hashable] struct {
	Key  T
	Rate float64
}

type Hitter[T pbtk.Hashable] interface {
	Add(k T) error
	Hits() []Hit[T]
	AppendHits(dst []Hit[T]) []Hit[T]
	Reset()
}
