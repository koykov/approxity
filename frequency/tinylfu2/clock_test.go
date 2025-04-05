package tinylfu

import "time"

type testClock struct {
	now time.Time
}

func (c *testClock) Now() time.Time {
	return c.now
}

func (c *testClock) UNow() uint64 {
	return uint64(c.now.UnixNano())
}

func (c *testClock) UNow32() uint32 {
	return uint32(c.now.Unix())
}

func (c *testClock) set(now time.Time) {
	c.now = now
}

func (c *testClock) add(d time.Duration) {
	c.now = c.now.Add(d)
}

func newTestClock(now time.Time) *testClock {
	return &testClock{now}
}
