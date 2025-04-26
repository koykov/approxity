package lsh

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
)

type tuple struct {
	ID         int     `tsv:"pair_ID"`
	A          []byte  `tsv:"sentence_A"`
	B          []byte  `tsv:"sentence_B"`
	RelScore   float64 `tsv:"relatedness_score"`
	Entailment []byte  `tsv:"entailment_judgment"`
}

type dataset struct {
	name   string
	tuples []tuple
}

const (
	datasetTrial = "https://raw.githubusercontent.com/brmson/dataset-sts/refs/heads/master/data/sts/sick2014/SICK_trial.txt"
	datasetTrain = "https://raw.githubusercontent.com/brmson/dataset-sts/refs/heads/master/data/sts/sick2014/SICK_train.txt"
	datasetTest  = "https://raw.githubusercontent.com/brmson/dataset-sts/refs/heads/master/data/sts/sick2014/SICK_test_annotated.txt"
)

var datasets []dataset

func init() {
	fread := func(rempath string) (dataset, error) {
		var ds dataset

		pos := strings.LastIndex(rempath, "/")
		if pos == -1 {
			return ds, fmt.Errorf("invalid path: %s", rempath)
		}
		fname := rempath[pos+1:]
		ds.name = strings.ReplaceAll(fname, ".txt", "")
		locpath := "/tmp/" + fname

		var r io.Reader

		if _, err := os.Stat(locpath); errors.Is(err, os.ErrNotExist) {
			u, err := http.Get(rempath)
			if err != nil {
				return ds, err
			}

			t, err := os.OpenFile(locpath, os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return ds, err
			}
			defer func() { _ = t.Close() }()
			_, err = io.Copy(t, u.Body)
			if err != nil {
				return ds, err
			}
			_ = u.Body.Close()
		}

		f, err := os.Open(locpath)
		if err != nil {
			return ds, err
		}
		defer func() { _ = f.Close() }()
		r = f

		rdr := csv.NewReader(r)
		rdr.Comma = '\t'
		for i := 0; ; i++ {
			rec, err := rdr.Read()
			if err == io.EOF {
				break
			}
			if i == 0 {
				continue
			}
			id, _ := strconv.Atoi(rec[0])
			score, _ := strconv.ParseFloat(rec[3], 64)
			t := tuple{
				ID:         id,
				A:          []byte(rec[1]),
				B:          []byte(rec[2]),
				RelScore:   score,
				Entailment: []byte(rec[4]),
			}
			ds.tuples = append(ds.tuples, t)
			if err != nil {
				break
			}
		}
		return ds, nil

	}
	if t, err := fread(datasetTrial); err == nil {
		datasets = append(datasets, t)
	}
	if t, err := fread(datasetTrain); err == nil {
		datasets = append(datasets, t)
	}
	if t, err := fread(datasetTest); err == nil {
		datasets = append(datasets, t)
	}
}

func TestMe[T []byte](t *testing.T, hash Hasher[T], distFn func([]uint64, []uint64, uint64) float64, numHashes uint64, expectAvgDist float64) {
	for i := 0; i < len(datasets); i++ {
		ds := &datasets[i]
		t.Run(ds.name, func(t *testing.T) {
			var s, c float64
			for j := 0; j < len(ds.tuples); j++ {
				tp := &ds.tuples[j]

				hash.Reset()
				_ = hash.Add(tp.A)
				h0 := hash.Hash()

				hash.Reset()
				_ = hash.Add(tp.B)
				h1 := hash.Hash()

				dist := (64 - float64(distFn(h0, h1, numHashes))) / 64
				s += dist
				c++
			}
			if avg := s / c; avg > expectAvgDist {
				t.Errorf("avg dist = %f, expected %f", avg, expectAvgDist)
			}
		})
	}
}

func BenchmarkMe[T []byte](b *testing.B, hash Hasher[T]) {
	stages := [][]byte{
		[]byte("foo"),
		[]byte("foobar"),
		[]byte("hello world"),
		[]byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit."),
		[]byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit. Mauris varius nisi erat, ac vulputate elit malesuada ut."),
		[]byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit. Mauris varius nisi erat, ac vulputate elit malesuada ut. Nulla facilisi. Vestibulum nec sapien nisl. Curabitur at elit fringilla, consectetur dui nec, maximus quam. Proin dui ipsum, venenatis nec est non, consectetur semper leo. Curabitur quis arcu ornare, malesuada nibh vel, maximus neque."),
	}
	for _, st := range stages {
		b.Run(fmt.Sprintf("add/%d", len(st)), func(b *testing.B) {
			b.SetBytes(int64(len(st)))
			b.ReportAllocs()
			b.ResetTimer()
			for j := 0; j < b.N; j++ {
				hash.Reset()
				_ = hash.Add(st)
			}
		})
	}
	for _, st := range stages {
		b.Run(fmt.Sprintf("hash/%d", len(st)), func(b *testing.B) {
			var buf []uint64
			b.SetBytes(int64(len(st)))
			b.ReportAllocs()
			b.ResetTimer()
			for j := 0; j < b.N; j++ {
				hash.Reset()
				_ = hash.Add(st)
				buf = hash.AppendHash(buf[:0])
				_ = buf
			}
		})
	}
}

func TestDistHamming(h0, h1 []uint64, _ uint64) (r float64) {
	bits := h0[0] ^ h1[0]
	for i := 0; i < 32; i++ {
		bs := bits & (1 << i)
		if bs != 0 {
			r += 1
		}
	}
	return r
}

func TestDistJaccard(h0, h1 []uint64, n uint64) (r float64) {
	if len(h1) < len(h0) {
		h0, h1 = h1, h0
	}
	for i := 0; i < len(h0); i++ {
		if h0[i] == h1[i] {
			r += 1
		}
	}
	if n == 0 {
		n = 1
	}
	return r / float64(n)
}
