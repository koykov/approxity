package hyperloglog

import (
	"encoding/binary"
	"math"
	"os"
)

// empirical bias correction pairs
// loads from local binary due to huge size
var bias [15]map[uint64]uint64

func biasfn(p uint64, e float64) float64 {
	_ = bias[14]
	v, ok := bias[p][math.Float64bits(e)]
	if !ok {
		return e
	}
	return math.Float64frombits(v)
}

func init() {
	fh, err := os.OpenFile("bias.bin", os.O_RDONLY, 0644)
	if err != nil {
		return
	}
	defer func() { _ = fh.Close() }()
	for i := 0; i < 15; i++ {
		var buf [8]byte
		if _, err = fh.Read(buf[:]); err != nil {
			return
		}
		n := binary.LittleEndian.Uint64(buf[:])
		bias[i] = make(map[uint64]uint64, n)
		for j := uint64(0); j < n; j++ {
			if _, err = fh.Read(buf[:]); err != nil {
				return
			}
			k := binary.LittleEndian.Uint64(buf[:])
			if _, err = fh.Read(buf[:]); err != nil {
				return
			}
			v := binary.LittleEndian.Uint64(buf[:])
			bias[i][k] = v
		}
	}
}
