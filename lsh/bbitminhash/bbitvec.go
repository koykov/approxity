package bbitminhash

// compact vector of b bits of values
type bbitvec struct {
	buf []uint64 // storage
	b   uint64   // number of lower bits
	off uint64   // offset of total used bits in buf
	c   uint64   // number of elements in buf
}

func newBbitvec(b uint64) *bbitvec {
	return &bbitvec{b: b}
}

func (v *bbitvec) add(val uint64) {
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

		// 2. Оставшиеся биты записываем в следующий элемент
		// - rest of bits moves to the next element
		next := lo >> remainings
		v.buf = append(v.buf, next)
	}

	v.off += v.b
	v.c++
}

func (v *bbitvec) each(fn func(uint64)) {
	if v.c == 0 {
		return
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

			// Обновляем позиции
			bpos += readBits
			bleft -= readBits

			if bpos%64 == 0 {
				idx++
			}
		}

		fn(curr)
	}
}

func (v *bbitvec) reset() {
	v.buf = v.buf[:0]
	v.off = 0
	v.c = 0
}
