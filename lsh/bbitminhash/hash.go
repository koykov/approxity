package bbitminhash

import (
	"github.com/koykov/byteseq"
	"github.com/koykov/pbtk/lsh"
	"github.com/koykov/pbtk/lsh/minhash"
)

func NewHasher[T byteseq.Q](conf *Config[T]) (lsh.Hasher[T], error) {
	if conf.B == 0 {
		return nil, lsh.ErrZeroB
	}
	if conf.B >= 64 {
		return nil, lsh.ErrBigB
	}
	if conf.Vector == nil {
		conf.Vector = newVector(conf.B)
	}
	h, err := minhash.NewHasher[T](&conf.Config)
	if err != nil {
		return nil, err
	}
	return h, nil
}
