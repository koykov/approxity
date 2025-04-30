package simtest

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Tuple struct {
	ID         int     `tsv:"pair_ID"`
	A          []byte  `tsv:"sentence_A"`
	B          []byte  `tsv:"sentence_B"`
	RelScore   float64 `tsv:"relatedness_score"`
	Entailment []byte  `tsv:"entailment_judgment"`
}

type Dataset struct {
	Name   string
	Tuples []Tuple
}

const (
	datasetTrial = "https://github.com/koykov/dataset/raw/refs/heads/master/similarity/SICK_trial.tsv"
	datasetTrain = "https://github.com/koykov/dataset/raw/refs/heads/master/similarity/SICK_train.tsv"
	datasetTest  = "https://github.com/koykov/dataset/raw/refs/heads/master/similarity/SICK_test_annotated.tsv"
)

var datasets []Dataset

func init() {
	fread := func(rempath string) (Dataset, error) {
		var ds Dataset

		pos := strings.LastIndex(rempath, "/")
		if pos == -1 {
			return ds, fmt.Errorf("invalid path: %s", rempath)
		}
		fname := rempath[pos+1:]
		ds.Name = strings.ReplaceAll(fname, ".tsv", "")
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
			t := Tuple{
				ID:         id,
				A:          []byte(rec[1]),
				B:          []byte(rec[2]),
				RelScore:   score,
				Entailment: []byte(rec[4]),
			}
			ds.Tuples = append(ds.Tuples, t)
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

func EachTestingDataset(f func(i int, ds *Dataset)) {
	for i := 0; i < len(datasets); i++ {
		f(i, &datasets[i])
	}
}
