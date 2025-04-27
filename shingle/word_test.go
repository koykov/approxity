package shingle

import (
	"strconv"
	"testing"
)

type wstage struct {
	name    string
	text    string
	tokens  map[uint][]string
	ctokens map[uint][]string
}

var wstages = []wstage{
	{
		name: "single word",
		text: "foobar!",
		tokens: map[uint][]string{
			2: {"foobar!"},
			3: {"foobar!"},
		},
		ctokens: map[uint][]string{
			2: {"foobar"},
			3: {"foobar"},
		},
	},
	{
		name: "phrase",
		text: "Stock markets hit record highs!",
		tokens: map[uint][]string{
			2: {"Stock markets", "markets hit", "hit record", "record highs!"},
			3: {"Stock markets hit", "markets hit record", "hit record highs!"},
			4: {"Stock markets hit record", "markets hit record highs!"},
		},
		ctokens: map[uint][]string{
			2: {"Stock markets", "markets hit", "hit record", "record highs"},
			3: {"Stock markets hit", "markets hit record", "hit record highs"},
			4: {"Stock markets hit record", "markets hit record highs"},
		},
	},
	{
		name: "sentence",
		text: "NASA's Mars rover discovers ancient riverbed - scientists thrilled!",
		tokens: map[uint][]string{
			2: {"NASA's Mars", "Mars rover", "rover discovers", "discovers ancient", "ancient riverbed", "riverbed -", "- scientists", "scientists thrilled!"},
			3: {"NASA's Mars rover", "Mars rover discovers", "rover discovers ancient", "discovers ancient riverbed", "ancient riverbed -", "riverbed - scientists", "- scientists thrilled!"},
			4: {"NASA's Mars rover discovers", "Mars rover discovers ancient", "rover discovers ancient riverbed", "discovers ancient riverbed -", "ancient riverbed - scientists", "riverbed - scientists thrilled!"},
		},
		ctokens: map[uint][]string{
			2: {"NASAs Mars", "Mars rover", "rover discovers", "discovers ancient", "ancient riverbed", "riverbed scientists", "scientists thrilled"},
			3: {"NASAs Mars rover", "Mars rover discovers", "rover discovers ancient", "discovers ancient riverbed", "ancient riverbed scientists", "riverbed scientists thrilled"},
			4: {"NASAs Mars rover discovers", "Mars rover discovers ancient", "rover discovers ancient riverbed", "discovers ancient riverbed scientists", "ancient riverbed scientists thrilled"},
		},
	},
	{
		name: "long sentence",
		text: "Tech giant Apple - after months of speculation - finally unveiled its revolutionary M4 AI chip at WWDC 2024, marking a major leap forward.",
		tokens: map[uint][]string{
			2: {"Tech giant", "giant Apple", "Apple -", "- after", "after months", "months of", "of speculation", "speculation -", "- finally", "finally unveiled", "unveiled its", "its revolutionary", "revolutionary M4", "M4 AI", "AI chip", "chip at", "at WWDC", "WWDC 2024,", "2024, marking", "marking a", "a major", "major leap", "leap forward."},
			3: {"Tech giant Apple", "giant Apple -", "Apple - after", "- after months", "after months of", "months of speculation", "of speculation -", "speculation - finally", "- finally unveiled", "finally unveiled its", "unveiled its revolutionary", "its revolutionary M4", "revolutionary M4 AI", "M4 AI chip", "AI chip at", "chip at WWDC", "at WWDC 2024,", "WWDC 2024, marking", "2024, marking a", "marking a major", "a major leap", "major leap forward."},
			4: {"Tech giant Apple -", "giant Apple - after", "Apple - after months", "- after months of", "after months of speculation", "months of speculation -", "of speculation - finally", "speculation - finally unveiled", "- finally unveiled its", "finally unveiled its revolutionary", "unveiled its revolutionary M4", "its revolutionary M4 AI", "revolutionary M4 AI chip", "M4 AI chip at", "AI chip at WWDC", "chip at WWDC 2024,", "at WWDC 2024, marking", "WWDC 2024, marking a", "2024, marking a major", "marking a major leap", "a major leap forward."},
		},
		ctokens: map[uint][]string{
			2: {"Tech giant", "giant Apple", "Apple after", "after months", "months of", "of speculation", "speculation finally", "finally unveiled", "unveiled its", "its revolutionary", "revolutionary M4", "M4 AI", "AI chip", "chip at", "at WWDC", "WWDC 2024", "2024 marking", "marking a", "a major", "major leap", "leap forward"},
			3: {"Tech giant Apple", "giant Apple after", "Apple after months", "after months of", "months of speculation", "of speculation finally", "speculation finally unveiled", "finally unveiled its", "unveiled its revolutionary", "its revolutionary M4", "revolutionary M4 AI", "M4 AI chip", "AI chip at", "chip at WWDC", "at WWDC 2024", "WWDC 2024 marking", "2024 marking a", "marking a major", "a major leap", "major leap forward"},
			4: {"Tech giant Apple after", "giant Apple after months", "Apple after months of", "after months of speculation", "months of speculation finally", "of speculation finally unveiled", "speculation finally unveiled its", "finally unveiled its revolutionary", "unveiled its revolutionary M4", "its revolutionary M4 AI", "revolutionary M4 AI chip", "M4 AI chip at", "AI chip at WWDC", "chip at WWDC 2024", "at WWDC 2024 marking", "WWDC 2024 marking a", "2024 marking a major", "marking a major leap", "a major leap forward"},
		},
	},
}

func TestWord(t *testing.T) {
	sheq := func(a, b []string) bool {
		if len(a) != len(b) {
			return false
		}
		for i := 0; i < len(a); i++ {
			if a[i] != b[i] {
				return false
			}
		}
		return true
	}
	for i := 0; i < len(wstages); i++ {
		st := &wstages[i]
		t.Run(st.name, func(t *testing.T) {
			t.Run("origin", func(t *testing.T) {
				for k, list := range st.tokens {
					t.Run(strconv.Itoa(int(k)), func(t *testing.T) {
						sh := NewWord[string](uint64(k), "")
						r := sh.Shingle(st.text)
						if !sheq(r, list) {
							t.Errorf("expected %+v, got %+v", list, r)
						}
					})
				}
			})
			t.Run("clean", func(t *testing.T) {
				for k, list := range st.ctokens {
					t.Run(strconv.Itoa(int(k)), func(t *testing.T) {
						sh := NewWord[string](uint64(k), CleanSetAll)
						r := sh.Shingle(st.text)
						if !sheq(r, list) {
							t.Errorf("expected %v, got %v", list, r)
						}
					})
				}
			})
		})
	}
}

func BenchmarkWord(b *testing.B) {
	for i := 0; i < len(wstages); i++ {
		st := &wstages[i]
		b.Run(st.name, func(b *testing.B) {
			b.Run("origin", func(b *testing.B) {
				for k, _ := range st.tokens {
					b.Run(strconv.Itoa(int(k)), func(b *testing.B) {
						sh := NewWord[string](uint64(k), "")
						var buf []string
						b.ReportAllocs()
						for j := 0; j < b.N; j++ {
							sh.Reset()
							buf = sh.AppendShingle(buf[:0], st.text)
						}
					})
				}
			})
			b.Run("clean", func(b *testing.B) {
				for k, _ := range st.ctokens {
					b.Run(strconv.Itoa(int(k)), func(b *testing.B) {
						sh := NewWord[string](uint64(k), CleanSetAll)
						var buf []string
						b.ReportAllocs()
						for j := 0; j < b.N; j++ {
							sh.Reset()
							buf = sh.AppendShingle(buf[:0], st.text)
						}
					})
				}
			})
		})
	}
}
