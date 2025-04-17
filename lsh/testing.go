package lsh

import (
	"encoding/csv"
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/koykov/pbtk"
)

type tuple struct {
	ID         int     `tsv:"pair_ID"`
	A          string  `tsv:"sentence_A"`
	B          string  `tsv:"sentence_B"`
	RelScore   float64 `tsv:"relatedness_score"`
	Entailment string  `tsv:"entailment_judgment"`
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
	fread := func(fpath string) (dataset, error) {
		var ds dataset

		f, err := http.Get(fpath)
		if err != nil {
			return ds, err
		}
		defer func() { _ = f.Body.Close() }()
		if pos := strings.LastIndex(fpath, "/"); pos != -1 {
			ds.name = strings.ReplaceAll(fpath[pos+1:], ".txt", "")
		}

		rdr := csv.NewReader(f.Body)
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
				A:          rec[1],
				B:          rec[2],
				RelScore:   score,
				Entailment: rec[4],
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

func TestMe[T pbtk.Hashable](t *testing.T, hash Hasher[T]) {
	for _, ds := range datasets {
		for _, t := range ds.tuples {
			_ = t
			// hash.Add(t)
		}
	}
}
