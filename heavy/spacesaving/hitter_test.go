package spacesaving

import (
	"testing"

	"github.com/koykov/hash/xxhash"
	"github.com/koykov/pbtk"
)

var testh = xxhash.Hasher64[[]byte]{}

func TestHitter(t *testing.T) {
	pbtk.EachTestingDataset(func(_ int, ds *pbtk.TestingDataset[[]byte]) {
		t.Run(ds.Name, func(t *testing.T) {
			h, err := NewHitter[[]byte](NewConfig(5, testh).
				WithBuckets(1).
				WithEWMA(0.5))
			if err != nil {
				t.Fatal(err)
			}
			for i := 0; i < len(ds.All); i++ {
				if err = h.Add(ds.All[i]); err != nil {
					t.Fatal(err)
				}
				switch {
				case i%5 == 0:
					_ = h.Add(ds.All[5])
				case i%10 == 0:
					_ = h.Add(ds.All[10])
				case i%100 == 0:
					_ = h.Add(ds.All[100])
				case i%1000 == 0:
					_ = h.Add(ds.All[1000])
				case i%10000 == 0:
					_ = h.Add(ds.All[10000])
				case i%100000 == 0:
					_ = h.Add(ds.All[100000])
				}
			}
			hits := h.Hits()
			r := make(map[string]float64, len(hits))
			for i := 0; i < len(hits); i++ {
				r[string(hits[i].Key)] = hits[i].Rate
			}
			t.Logf("%#v", r)
		})
	})
}
