package quotient

type bucket uint64

type btype uint64

const (
	btypeOccupied btype = iota
	btypeContinuation
	btypeShifted
)
const xor7 = int64(^7)

func newBucket(r uint64) bucket {
	return bucket(int64(r<<3) & xor7)
}

func (b *bucket) empty() bool {
	return *b&7 == 0
}

func (b *bucket) setbit(bt btype) {
	switch bt {
	case btypeOccupied:
		*b |= 1
	case btypeContinuation:
		*b |= 2
	case btypeShifted:
		*b |= 4
	}
}

func (b *bucket) clearbit(bt btype) {
	var cb int64
	switch bt {
	case btypeOccupied:
		cb = ^1
	case btypeContinuation:
		cb = ^2
	case btypeShifted:
		cb = ^4
	}
	*b = *b & bucket(cb)
}

func (b *bucket) checkbit(bt btype) bool {
	return *b&(1<<bt) != 0
}

func (b *bucket) rem() uint64 {
	return uint64(*b >> 3)
}
func (b *bucket) raw() uint64 {
	return uint64(*b)
}
