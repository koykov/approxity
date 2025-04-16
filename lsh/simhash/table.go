package simhash

var (
	btable = [2]int64{-1, 1}
	rtable = [vecsz]uint64{}
)

func init() {
	for i := uint64(0); i < vecsz; i++ {
		rtable[i] = uint64(1) << i
	}
}
