package amq

import (
	"bufio"
	"context"
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type dataset[T []byte] struct {
	name     string
	positive []T
	negative []T
	all      []T
}

var datasets []dataset[[]byte]

func init() {
	fread := func(dst [][]byte, path string) ([][]byte, error) {
		fh, err := os.Open(path)
		if err != nil {
			return dst, err
		}
		defer func() { _ = fh.Close() }()
		scr := bufio.NewScanner(fh)
		for scr.Scan() {
			if b := scr.Bytes(); len(b) > 0 {
				dst = append(dst, append([]byte(nil), b...))
			}
		}
		return dst, scr.Err()
	}
	probes := []string{
		"testdata",
		"../testdata",
	}
	for _, path := range probes {
		_ = filepath.Walk(path, func(cpath string, info fs.FileInfo, err error) error {
			if info == nil || !info.IsDir() || cpath == path {
				return nil
			}
			var ds dataset[[]byte]
			if ds.positive, err = fread(ds.positive, cpath+"/positive.txt"); err != nil {
				return err
			}
			if ds.negative, err = fread(ds.negative, cpath+"/negative.txt"); err != nil {
				return err
			}
			ds.name, ds.all = info.Name(), append(ds.positive, ds.negative...)
			datasets = append(datasets, ds)
			return nil
		})
	}
	// Try to compose dataset based on system's dictionaries.
	sysDS := dataset[[]byte]{name: "system/dict"}
	if words, err := fread(nil, "/usr/share/dict/words"); err == nil && len(words) > 0 {
		sysDS.all = words
		for i := 0; i < len(words); i++ {
			if i%2 == 0 {
				sysDS.positive = append(sysDS.positive, words[i])
			} else {
				sysDS.negative = append(sysDS.negative, words[i])
			}
		}
		datasets = append(datasets, sysDS)
	}
}

func TestMe[T []byte](t *testing.T, f Filter[T]) {
	for i := 0; i < len(datasets); i++ {
		ds := &datasets[i]
		t.Run(ds.name, func(t *testing.T) {
			f.Reset()
			for j := 0; j < len(ds.positive); j++ {
				_ = f.Set(ds.positive[j])
			}
			var falseNegative, falsePositive int
			for j := 0; j < len(ds.negative); j++ {
				if f.Contains(ds.negative[j]) {
					falsePositive++
				}
			}
			if falsePositive > 0 {
				// Just warn, it's OK to have small amount of false positives.
				t.Logf("%d of %d negatives (%d total) gives false positive value", falsePositive, len(ds.negative), len(ds.all))
			}
			for j := 0; j < len(ds.positive); j++ {
				if !f.Contains(ds.positive[j]) {
					falseNegative++
				}
			}
			if falseNegative > 0 {
				t.Errorf("%d of %d positives (%d total) gives false negative value", falseNegative, len(ds.positive), len(ds.all))
			}
		})
	}
}

func TestMeConcurrently[T []byte](t *testing.T, f Filter[T]) {
	for i := 0; i < len(datasets); i++ {
		ds := &datasets[i]
		t.Run(ds.name, func(t *testing.T) {
			f.Reset()
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			var wg sync.WaitGroup
			wg.Add(3)

			go func() {
				defer wg.Done()
				for j := 0; ; j++ {
					select {
					case <-ctx.Done():
						return
					default:
						_ = f.Set(ds.positive[j%len(ds.positive)])
					}
				}
			}()

			go func() {
				defer wg.Done()
				for j := 0; ; j++ {
					select {
					case <-ctx.Done():
						return
					default:
						_ = f.Unset(ds.all[j%len(ds.all)])
					}
				}
			}()

			go func() {
				defer wg.Done()
				for j := 0; ; j++ {
					select {
					case <-ctx.Done():
						return
					default:
						f.Contains(ds.all[(j % len(ds.all))])
					}
				}
			}()

			wg.Wait()
		})
	}
}

func BenchMe[T []byte](b *testing.B, f Filter[T]) {
	for i := 0; i < len(datasets); i++ {
		ds := &datasets[i]
		b.Run(ds.name, func(b *testing.B) {
			f.Reset()
			for j := 0; j < len(ds.positive); j++ {
				_ = f.Set(ds.positive[j])
			}
			b.ReportAllocs()
			b.ResetTimer()
			for k := 0; k < b.N; k++ {
				f.Contains(ds.all[k%len(ds.all)])
			}
		})
	}
}

func BenchMeConcurrently[T []byte](b *testing.B, f Filter[T]) {
	for i := 0; i < len(datasets); i++ {
		ds := &datasets[i]
		b.Run(ds.name, func(b *testing.B) {
			f.Reset()
			b.ReportAllocs()
			b.RunParallel(func(pb *testing.PB) {
				var j uint64 = math.MaxUint64
				for pb.Next() {
					ci := atomic.AddUint64(&j, 1)
					switch ci % 100 {
					case 99:
						_ = f.Set(ds.positive[ci%uint64(len(ds.positive))])
					case 98:
						_ = f.Unset(ds.all[ci%uint64(len(ds.all))])
					default:
						f.Contains(ds.all[ci%uint64(len(ds.all))])
					}
				}
			})
		})
	}
}
