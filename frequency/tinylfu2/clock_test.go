package tinylfu

import "time"

type testClock struct {
	now time.Time
}

func (c *testClock) Now() time.Time {
	return c.now
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
