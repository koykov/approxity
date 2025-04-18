package simhash

var (
	btable = [2]int64{-1, 1}
	rtable = [bucketsz]uint64{}
)

func init() {
	for i := uint64(0); i < bucketsz; i++ {
		rtable[i] = uint64(1) << i
	}
}
