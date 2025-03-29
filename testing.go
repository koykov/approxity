package pbtk

import (
	"bufio"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
)

type TestingDataset[T []byte] struct {
	Name      string
	Positives []T
	Negatives []T
	All       []T
}

var datasets []TestingDataset[[]byte]

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
		"../../testdata",
	}
	for _, path := range probes {
		_ = filepath.Walk(path, func(cpath string, info fs.FileInfo, err error) error {
			if info == nil || !info.IsDir() || cpath == path {
				return nil
			}
			var ds TestingDataset[[]byte]
			if ds.Positives, err = fread(ds.Positives, cpath+"/positive.txt"); err != nil {
				return err
			}
			if ds.Negatives, err = fread(ds.Negatives, cpath+"/negative.txt"); err != nil {
				return err
			}
			ds.Name, ds.All = info.Name(), append(ds.Positives, ds.Negatives...)
			datasets = append(datasets, ds)
			return nil
		})
	}
	// Try to compose dataset based on system's dictionaries.
	sysDS := TestingDataset[[]byte]{Name: "system/dict"}
	if words, err := fread(nil, "/usr/share/dict/words"); err == nil && len(words) > 0 {
		sysDS.All = words
		for i := 0; i < len(words); i++ {
			if i%2 == 0 {
				sysDS.Positives = append(sysDS.Positives, words[i])
			} else {
				sysDS.Negatives = append(sysDS.Negatives, words[i])
			}
		}
		datasets = append(datasets, sysDS)
	}
	{
		const (
			localDS  = "/tmp/probds_eng_voc.txt"
			remoteDS = "https://raw.githubusercontent.com/koykov/dataset/master/vocabulary/freelang/English.txt"
		)
		// Try to compose dataset based on remote English vocabulary.
		var (
			vocDS TestingDataset[[]byte]
			err   error
		)
		// - try local cache first
		if vocDS.All, err = fread(vocDS.All, localDS); err != nil {
			// - negative, try fetch remote data and save to local cache
			var resp *http.Response
			resp, err = http.Get(remoteDS)
			if err == nil && resp.StatusCode == http.StatusOK {
				defer func() { _ = resp.Body.Close() }()
				var contents []byte
				if contents, err = io.ReadAll(resp.Body); err == nil {
					err = os.WriteFile(localDS, contents, 0644)
				}
				// - try local cache again
				vocDS.All, err = fread(vocDS.All, localDS)
			}
		}
		if err == nil {
			// - fill up positives and negatives
			for i := 0; i < len(vocDS.All); i++ {
				if i%2 == 0 {
					vocDS.Positives = append(vocDS.Positives, vocDS.All[i])
				} else {
					vocDS.Negatives = append(vocDS.Negatives, vocDS.All[i])
				}
			}
			vocDS.Name = "english vocabulary"
			datasets = append(datasets, vocDS)
		}
	}
}

func EachTestingDataset(f func(i int, ds *TestingDataset[[]byte])) {
	for i := 0; i < len(datasets); i++ {
		f(i, &datasets[i])
	}
}
