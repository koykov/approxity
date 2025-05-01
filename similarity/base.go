package similarity

import (
	"github.com/koykov/byteseq"
	"github.com/koykov/pbtk/lsh"
)

type Base[T byteseq.Q] struct {
	buf []uint64
}

func (sim *Base[T]) VectorizePair(hasher lsh.Hasher[T], a, b T) ([]uint64, []uint64, error) {
	if err := hasher.Add(a); err != nil {
		return nil, nil, err
	}
	var mid int
	sim.buf = hasher.AppendHash(sim.buf[:0])
	mid = len(sim.buf)

	hasher.Reset()
	if err := hasher.Add(b); err != nil {
		return nil, nil, err
	}
	sim.buf = hasher.AppendHash(sim.buf)
	return sim.buf[:mid], sim.buf[mid:], nil
}

func (sim *Base[T]) Reset() {
	sim.buf = sim.buf[:0]
}
