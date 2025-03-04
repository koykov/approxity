package hyperloglog

import (
	"encoding/binary"
	"math"
	"os"
	"sort"
	"unsafe"
)

// empirical bias correction pairs
// loads from local binary due to huge size
var bias [15][][2]float64

func biasfn(p uint64, e float64) float64 {
	_ = bias[14]
	const k = 6 // K-nn
	var a [96]byte
	keys := *(*[k][2]float64)(unsafe.Pointer(&a))
	// var keys [k][2]float64
	ssize := len(bias[p])
	var eidx int
	{
		// lower_bound
		for i := 0; i < ssize; i++ {
			if bias[p][i][0] < e {
				eidx = i
				continue
			}
			break
		}
	}
	{
		// partial_sort_copy
		lo, hi := 0, ssize
		if eidx > k {
			lo = eidx - k
		}
		if eidx+k < ssize {
			hi = eidx + k
		}
		for i := lo; i < hi; i++ {
			keys[i-lo][0], keys[i-lo][1] = bias[p][i][0], bias[p][i][1]
		}
		sort.Slice(keys[:], func(i, j int) bool {
			return math.Abs(keys[i][0]-e) < math.Abs(keys[j][0]-e)
		})
	}
	var s, ws float64
	{
		// accumulate
		for i := 0; i < k; i++ {
			s += keys[i][1] * 1 / (math.Abs(keys[i][0]-e) + 1e-5)
			ws += 1 / (math.Abs(keys[i][0]-e) + 1e-5)
		}
	}
	return s / ws
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
		for j := uint64(0); j < n; j++ {
			if _, err = fh.Read(buf[:]); err != nil {
				return
			}
			dist := binary.LittleEndian.Uint64(buf[:])
			if _, err = fh.Read(buf[:]); err != nil {
				return
			}
			bias_ := binary.LittleEndian.Uint64(buf[:])
			bias[i] = append(bias[i], [2]float64{math.Float64frombits(dist), math.Float64frombits(bias_)})
		}
	}
}
