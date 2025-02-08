package bloom

import "testing"

var dataset = []struct {
	pos, neg []string
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

func assertBool(tb testing.TB, value, expected bool) {
	if value != expected {
		tb.Errorf("expected %v, got %v", expected, value)
	}
}

func TestFilter(t *testing.T) {
	for i := 0; i < len(dataset); i++ {
		ds := &dataset[i]
		t.Run("", func(t *testing.T) {
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
	}
}

func BenchmarkFilter(b *testing.B) {
	for i := 0; i < len(dataset); i++ {
		ds := &dataset[i]
		b.Run("", func(b *testing.B) {
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
			all := make([]string, 0, len(ds.pos)+len(ds.neg))
			all = append(all, ds.pos...)
			all = append(all, ds.neg...)
			for k := 0; k < b.N; k++ {
				f.Check(&all[k%len(all)])
			}
		})
	}
}
