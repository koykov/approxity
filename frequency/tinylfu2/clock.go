package tinylfu

import "time"

type Clock interface {
	Now() time.Time
	UNow() uint64
	UNow32() uint32
}

type nativeClock struct{}

func (nativeClock) Now() time.Time {
	return time.Now()
}

func (nativeClock) UNow() uint64 {
	return uint64(time.Now().UnixNano())
}

func (nativeClock) UNow32() uint32 {
	return uint32(time.Now().Unix())
}
