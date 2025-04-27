package minhash

import (
	"math"
	"strconv"
	"sync"
	"unsafe"

	"github.com/koykov/byteseq"
	"github.com/koykov/openrt"
	"github.com/koykov/pbtk"
	"github.com/koykov/pbtk/lsh"
	"github.com/koykov/simd/memset64"
)

type hash[T byteseq.Q] struct {
	b      pbtk.Base[T]
	conf   *Config[T]
	vector []uint64
	token  []T
	buf    []byte
	once   sync.Once

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
	n := len(h.token)

	if cap(h.vector) < n {
		h.vector = make([]uint64, n)
	}
	h.vector = h.vector[:n]
	memset64.Memset(h.vector, math.MaxUint64)
	for i := 0; i < n; i++ {
		for j := uint64(0); j < h.conf.K; j++ {
			h.buf = append(h.buf[:0], h.token[i]...)
			h.buf = strconv.AppendUint(h.buf, uint64(j), 10)
			hsum := h.conf.Algo.Sum64(h.buf)
			h.vector[i] = min(h.vector[i], hsum)
		}
	}
	return nil
}

func (h *hash[T]) Hash() []uint64 {
	r := make([]uint64, 0, len(h.vector))
	return h.AppendHash(r)
}

func (h *hash[T]) AppendHash(dst []uint64) []uint64 {
	return append(dst, h.vector...)
}

func (h *hash[T]) Reset() {
	if len(h.vector) == 0 {
		return
	}
	h.token = h.token[:0]
	h.buf = h.buf[:0]
	h.conf.Shingler.Reset()
	openrt.MemclrUnsafe(unsafe.Pointer(&h.vector[0]), len(h.vector)*8)
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
}
