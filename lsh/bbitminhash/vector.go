package bbitminhash

// compact vector of b bits of values
type vector struct {
	buf []uint64 // storage
	b   uint64   // number of lower bits
	off uint64   // offset of total used bits in buf
	c   uint64   // number of elements in buf
}

func newVector(b uint64) *vector {
	return &vector{b: b}
}

func (v *vector) Grow(cap_ uint64) {
	total := cap_ * v.b
	sz := (total + 63) / 64

	if uint64(cap(v.buf)) < sz {
		nbuf := make([]uint64, sz)
		copy(nbuf, v.buf)
		v.buf = nbuf
	}
	v.buf = v.buf[:sz]
	v.c = cap_
}

func (v *vector) Add(val uint64) {
	lo := val & ((1 << v.b) - 1)

	// position in buf
	idx := int(v.off / 64) // index in buf
	boff := v.off % 64     // offset in current element

	if idx >= len(v.buf) {
		v.buf = append(v.buf, 0)
	}

	if boff+v.b <= 64 {
		// happy path - lo bits hits into free space in current element
		v.buf[idx] |= lo << boff
	} else {
		// split lo into two parts:
		// - remaining bits that hits into free space in current element
		remainings := 64 - boff
		v.buf[idx] |= lo << boff

		// - rest of bits moves to the next element
		next := lo >> remainings
		v.buf = append(v.buf, next)
	}

	v.off += v.b
	v.c++
}

func (v *vector) SetMin(pos, val uint64) {
	curr := v.Get(pos)
	reduced := v.reduce(val)
	if reduced < curr {
		v.set(pos, reduced)
	}
}

func (v *vector) Get(pos uint64) uint64 {
	if pos >= v.c {
		return 0
	}

	bpos := pos * v.b
	idx := bpos / 64
	boff := bpos % 64

	if boff+v.b <= 64 {
		return (v.buf[idx] >> boff) & ((1 << v.b) - 1)
	}

	rem := 64 - boff
	lo := v.buf[idx] >> boff
	hi := v.buf[idx+1] & ((1 << (v.b - rem)) - 1)
	return lo | (hi << rem)
}

func (v *vector) Memset(val uint64) {
	if v.c == 0 {
		return
	}

	reduced := v.reduce(val)
	for i := uint64(0); i < v.c; i++ {
		v.set(i, reduced)
	}
}

func (v *vector) set(pos uint64, val uint64) {
	bpos := pos * v.b
	idx := bpos / 64
	boff := bpos % 64
	mask := uint64(((1 << v.b) - 1) << boff)

	// clear current bits
	v.buf[idx] &^= mask
	// write new bits
	v.buf[idx] |= (val & ((1 << v.b) - 1)) << boff

	if boff+v.b > 64 {
		rem := v.b - (64 - boff)
		next := idx + 1
		if next >= uint64(len(v.buf)) {
			v.buf = append(v.buf, 0)
		}

		v.buf[next] &^= (1 << rem) - 1
		v.buf[next] |= val >> (64 - boff)
	}
}

func (v *vector) AppendAll(dst []uint64) []uint64 {
	if v.c == 0 {
		return dst
	}

	var (
		bpos  uint64 // current position in bits since start of buf
		idx   uint64 // current index in buf
		bleft = v.b  // number of bits to read
		curr  uint64 // current accumulated value
	)

	for i := uint64(0); i < v.c; i++ {
		bleft = v.b
		curr = 0

		for bleft > 0 {
			// how many bits can be read from current element
			bitsInElem := 64 - (bpos % 64)
			readBits := bitsInElem
			if bleft < readBits {
				readBits = bleft
			}

			mask := (uint64(1) << readBits) - 1
			shift := bpos % 64
			value := (v.buf[idx] >> shift) & mask
			curr |= value << (v.b - bleft)

			bpos += readBits
			bleft -= readBits

			if bpos%64 == 0 {
				idx++
			}
		}

		dst = append(dst, curr)
	}
	return dst
}

func (v *vector) Len() uint64 {
	return v.c
}

func (v *vector) Reset() {
	v.buf = v.buf[:0]
	v.off = 0
	v.c = 0
}

func (v *vector) reduce(val uint64) (r uint64) {
	r = val & ((1 << v.b) - 1)
	r = r | (r << v.b)
	r = r | (r << (2 * v.b))
	return
}
