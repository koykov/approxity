package tinylfu

import (
	"io"
	"testing"
	"time"

	"github.com/koykov/hash/xxhash"
	"github.com/koykov/pbtk/frequency"
)

const (
	testConfidence = 0.99999
	testEpsilon    = 0.00001
)

var testh = xxhash.Hasher64[[]byte]{}

func TestEstimator(t *testing.T) {
	t.Run("dataset", func(t *testing.T) {
		est, err := NewEstimator[[]byte](NewConfig(testConfidence, testEpsilon, testh))
		if err != nil {
			t.Fatal(err)
		}
		frequency.TestMe(t, frequency.NewTestAdapter(est))
	})
	t.Run("decay", func(t *testing.T) {
		tryclose := func(est frequency.Estimator[string]) error {
			if c, ok := any(est).(io.Closer); ok {
				return c.Close()
			}
			return nil
		}
		t.Run("counter", func(t *testing.T) {
			est, _ := NewEstimator[string](NewConfig(0.99, 0.01, testh).
				WithDecayLimit(20))
			for i := 0; i < 10; i++ {
				_ = est.Add("foobar")
				_ = est.Add("qwerty")
			}
			_ = est.Add("final")
			e0, e1 := est.Estimate("foobar"), est.Estimate("qwerty")
			_ = tryclose(est)
			if e0 != 5 || e1 != 5 {
				t.Fatalf("unexpected estimates: %d, %d", e0, e1)
			}
		})
		t.Run("time interval", func(t *testing.T) {
			est, _ := NewEstimator[string](NewConfig(0.99, 0.01, testh).
				WithDecayInterval(time.Millisecond * 100))
			for i := 0; i < 10; i++ {
				_ = est.Add("foobar")
				_ = est.Add("qwerty")
			}
			time.Sleep(time.Millisecond * 110)
			e0, e1 := est.Estimate("foobar"), est.Estimate("qwerty")
			_ = tryclose(est)
			if e0 != 5 || e1 != 5 {
				t.Fatalf("unexpected estimates: %d, %d", e0, e1)
			}
		})
		t.Run("force", func(t *testing.T) {
			fd := testForceDecay{}
			est, _ := NewEstimator[string](NewConfig(0.99, 0.01, testh).
				WithForceDecayNotifier(&fd))
			for i := 0; i < 10; i++ {
				_ = est.Add("foobar")
				_ = est.Add("qwerty")
				if i == 3 || i == 6 || i == 9 {
					fd.trigger()
				}
			}
			e0, e1 := est.Estimate("foobar"), est.Estimate("qwerty")
			_ = tryclose(est)
			if e0 != 2 || e1 != 2 {
				t.Fatalf("unexpected estimates: %d, %d", e0, e1)
			}
		})
		t.Run("mixed", func(t *testing.T) {

		})
	})
}

func BenchmarkEstimator(b *testing.B) {
	b.Run("dataset", func(b *testing.B) {
		est, err := NewEstimator[[]byte](NewConfig(testConfidence, testEpsilon, testh))
		if err != nil {
			b.Fatal(err)
		}
		frequency.BenchMe(b, frequency.NewTestAdapter(est))
	})
	b.Run("dataset parallel", func(b *testing.B) {
		est, err := NewEstimator[[]byte](NewConfig(testConfidence, testEpsilon, testh).
			WithConcurrency().WithWriteAttemptsLimit(5))
		if err != nil {
			b.Fatal(err)
		}
		frequency.BenchMeConcurrently(b, frequency.NewTestAdapter(est))
	})
}

type testForceDecay struct {
	c chan struct{}
}

func (d *testForceDecay) Notify() <-chan struct{} {
	return d.c
}

func (d *testForceDecay) trigger() {
	if d.c == nil {
		d.c = make(chan struct{})
	}
	d.c <- struct{}{}
}

func (d *testForceDecay) close() {
	close(d.c)
}
