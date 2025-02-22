package cuckoo

import (
	"testing"

	"github.com/koykov/hash/metro"
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
			f, err := NewFilter(&Config{
				Size:   1e6,
				Hasher: metro.Hasher64[[]byte]{Seed: 1234},
			})
			if err != nil {
				t.Fatal(err)
			}
			for j := 0; j < len(ds.pos); j++ {
				_ = f.Set(ds.pos[j])
			}
			for j := 0; j < len(ds.neg); j++ {
				assertBool(t, f.Contains(ds.neg[j]), false)
			}
			for j := 0; j < len(ds.neg); j++ {
				assertBool(t, f.Contains(ds.pos[j]), true)
			}
			f.Reset()
			for j := 0; j < len(ds.pos); j++ {
				assertBool(t, f.Contains(ds.pos[j]), false)
			}
		})
	}
}

func BenchmarkFilter(b *testing.B) {
	for i := 0; i < len(dataset); i++ {
		ds := &dataset[i]
		b.Run("sync", func(b *testing.B) {
			f, err := NewFilter(&Config{
				Size:   1e6,
				Hasher: metro.Hasher64[[]byte]{Seed: 1234},
			})
			if err != nil {
				b.Fatal(err)
			}
			for j := 0; j < len(ds.pos); j++ {
				_ = f.Set(ds.pos[j])
			}
			b.ReportAllocs()
			b.ResetTimer()
			for k := 0; k < b.N; k++ {
				f.Contains(&ds.all[k%len(ds.all)])
			}
		})
	}
}
