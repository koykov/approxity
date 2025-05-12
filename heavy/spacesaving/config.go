package spacesaving

import (
	"github.com/koykov/pbtk"
	"github.com/koykov/pbtk/heavy"
)

type Config struct {
	K             uint64
	Hasher        pbtk.Hasher
	Buckets       uint64
	MetricsWriter heavy.MetricsWriter
}
