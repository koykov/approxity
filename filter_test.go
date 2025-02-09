package bloom

import (
	"context"
	"math"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

var dataset = []struct {
	pos, neg, all []string
}{
	{
		pos: []string{
			"abound", "abounds", "abundance", "abundant", "accessible",
			"bloom", "blossom", "bolster", "bonny", "bonus", "bonuses",
			"coherent", "cohesive", "colorful", "comely", "comfort",
			"gems", "generosity", "generous", "generously", "genial"},
		neg: []string{
			"bluff", "cheater", "hate", "war", "humanity",
			"racism", "hurt", "nuke", "gloomy", "facebook",
			"twitter", "google", "youtube", "comedy"},
	},
}

func init() {
	for i := 0; i < len(dataset); i++ {
		ds := &dataset[i]
		ds.all = make([]string, 0, len(ds.pos)+len(ds.neg))
		ds.all = append(ds.all, ds.pos...)
		ds.all = append(ds.all, ds.neg...)
	}
}

func assertBool(tb testing.TB, value, expected bool) {
	if value != expected {
		tb.Errorf("expected %v, got %v", expected, value)
	}
}

func TestFilter(t *testing.T) {
	for i := 0; i < len(dataset); i++ {
		ds := &dataset[i]
		t.Run("sync", func(t *testing.T) {
			f, err := NewFilter(NewConfig(1e5, &hasherStringCRC64{}).
				WithHashCheckLimit(3))
			if err != nil {
				t.Fatal(err)
			}
			for j := 0; j < len(ds.pos); j++ {
				_ = f.Set(ds.pos[j])
			}
			for j := 0; j < len(ds.neg); j++ {
				assertBool(t, f.Check(ds.neg[j]), false)
			}
			for j := 0; j < len(ds.neg); j++ {
				assertBool(t, f.Check(ds.pos[j]), true)
			}
		})
		t.Run("concurrent", func(t *testing.T) {
			f, err := NewFilter(NewConfig(1e5, &hasherStringCRC64{}).
				WithHashCheckLimit(3).
				WithConcurrency().WithWriteAttemptsLimit(5))
			if err != nil {
				t.Fatal(err)
			}

			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			var wg sync.WaitGroup
			wg.Add(3)

			go func() {
				defer wg.Done()
				for i := 0; ; i++ {
					select {
					case <-ctx.Done():
						return
					default:
						_ = f.Set(&ds.pos[i%len(ds.pos)])
					}
				}
			}()

			go func() {
				defer wg.Done()
				for {
					select {
					case <-ctx.Done():
						return
					default:
						_ = f.Clear(&ds.all[i%len(ds.all)])
					}
				}
			}()

			go func() {
				defer wg.Done()
				for {
					select {
					case <-ctx.Done():
						return
					default:
						f.Check(&ds.all[(i % len(ds.all))])
					}
				}
			}()

			wg.Wait()
		})
	}
}

func BenchmarkFilter(b *testing.B) {
	for i := 0; i < len(dataset); i++ {
		ds := &dataset[i]
		b.Run("sync", func(b *testing.B) {
			f, err := NewFilter(NewConfig(1e5, &hasherStringCRC64{}).
				WithHashCheckLimit(3))
			if err != nil {
				b.Fatal(err)
			}
			for j := 0; j < len(ds.pos); j++ {
				_ = f.Set(ds.pos[j])
			}
			b.ReportAllocs()
			b.ResetTimer()
			for k := 0; k < b.N; k++ {
				f.Check(&ds.all[k%len(ds.all)])
			}
		})
		b.Run("concurrent", func(b *testing.B) {
			b.ReportAllocs()

			f, _ := NewFilter(NewConfig(1e5, &hasherStringCRC64{}).
				WithHashCheckLimit(3).
				WithConcurrency().WithWriteAttemptsLimit(5))

			b.RunParallel(func(pb *testing.PB) {
				var i uint64 = math.MaxUint64
				for pb.Next() {
					ci := atomic.AddUint64(&i, 1)
					switch ci % 5 {
					case 4:
						_ = f.Set(&ds.pos[ci%uint64(len(ds.pos))])
					case 3:
						_ = f.Clear(&ds.all[ci%uint64(len(ds.all))])
					default:
						f.Check(&ds.all[ci%uint64(len(ds.all))])
					}
				}
			})
		})
	}
}
