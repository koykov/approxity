package cosine

import (
	"github.com/koykov/byteseq"
	"github.com/koykov/pbtk/lsh"
)

type Config[T byteseq.Q] struct {
	LSH lsh.Hasher[T]
}
