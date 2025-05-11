package heavy

import (
	"github.com/koykov/pbtk"
)

type Hitter[T pbtk.Hashable] interface {
	Add(k T) error
	Hits() []T
	AppendHits(dst []T) []T
	Reset()
}
