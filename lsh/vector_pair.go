package lsh

import "github.com/koykov/byteseq"

// VectorPair is a base type to build vectors of a pair of strings.
type VectorPair[T byteseq.Q] struct {
	buf []uint64
}

func (vp *VectorPair[T]) VectorizePair(hasher Hasher[T], a, b T) ([]uint64, []uint64, error) {
	if err := hasher.Add(a); err != nil {
		return nil, nil, err
	}
	var mid int
	vp.buf = hasher.AppendHash(vp.buf[:0])
	mid = len(vp.buf)

	hasher.Reset()
	if err := hasher.Add(b); err != nil {
		return nil, nil, err
	}
	vp.buf = hasher.AppendHash(vp.buf)
	return vp.buf[:mid], vp.buf[mid:], nil
}

func (vp *VectorPair[T]) Reset() {
	vp.buf = vp.buf[:0]
}
