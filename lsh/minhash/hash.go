package minhash

import (
	"math"
	"strconv"
	"sync"

	"github.com/koykov/byteseq"
	"github.com/koykov/pbtk"
	"github.com/koykov/pbtk/lsh"
)

type hash[T byteseq.Q] struct {
	b     pbtk.Base[T]
	conf  *Config[T]
	token []T
	buf   []byte
	once  sync.Once

	err error
}

func NewHasher[T byteseq.Q](conf *Config[T]) (lsh.Hasher[T], error) {
	h := &hash[T]{conf: conf.copy()}
	if h.once.Do(h.init); h.err != nil {
		return nil, h.err
	}
	return h, nil
}

func (h *hash[T]) Add(value T) error {
	if h.once.Do(h.init); h.err != nil {
		return h.err
	}
	h.token = h.conf.Shingler.AppendShingle(h.token, value)
	n := uint64(len(h.token))

	h.vec().Grow(n)
	h.vec().Memset(math.MaxUint64)
	for i := uint64(0); i < n; i++ {
		for j := uint64(0); j < h.conf.K; j++ {
			h.buf = append(h.buf[:0], h.token[i]...)
			h.buf = strconv.AppendUint(h.buf, j, 10)
			hsum := h.conf.Algo.Sum64(h.buf)
			h.vec().SetMin(i, hsum)
		}
	}
	return nil
}

func (h *hash[T]) Hash() []uint64 {
	r := make([]uint64, 0, h.vec().Len())
	return h.AppendHash(r)
}

func (h *hash[T]) AppendHash(dst []uint64) []uint64 {
	return h.vec().AppendAll(dst)
}

func (h *hash[T]) Reset() {
	if h.vec().Len() == 0 {
		return
	}
	h.token = h.token[:0]
	h.buf = h.buf[:0]
	h.conf.Shingler.Reset()
	h.vec().Reset()
}

func (h *hash[T]) init() {
	if h.conf.Algo == nil {
		h.err = pbtk.ErrNoHasher
		return
	}
	if h.conf.K == 0 {
		h.err = lsh.ErrZeroK
		return
	}
	if h.conf.Shingler == nil {
		h.err = lsh.ErrNoShingler
		return
	}
	if h.conf.Vector == nil {
		h.conf.Vector = &DefaultVector{}
	}
}

func (h *hash[T]) vec() Vector {
	return h.conf.Vector
}
