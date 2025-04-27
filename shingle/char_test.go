package shingle

import (
	"strconv"
	"testing"
)

var cstages = []stage{
	{
		name: "single word",
		text: "Hello!",
		tokens: map[uint][]string{
			2: {"He", "el", "ll", "lo", "o!"},
			3: {"Hel", "ell", "llo", "lo!"},
			4: {"Hell", "ello", "llo!"},
		},
		ctokens: map[uint][]string{
			2: {"He", "el", "ll", "lo"},
			3: {"Hel", "ell", "llo"},
			4: {"Hell", "ello"},
		},
	},
	{
		name: "phrase",
		text: "@user: $100 😊",
		tokens: map[uint][]string{
			2: {"@u", "us", "se", "er", "r:", ": ", " $", "$1", "10", "00", "0 ", " 😊"},
			3: {"@us", "use", "ser", "er:", "r: ", ": $", " $1", "$10", "100", "00 ", "0 😊"},
			4: {"@use", "user", "ser:", "er: ", "r: $", ": $1", " $10", "$100", "100 ", "00 😊"},
		},
		ctokens: map[uint][]string{
			2: {"us", "se", "er", "r ", " 1", "10", "00", "0 ", " 😊"},
			3: {"use", "ser", "er ", "r 1", " 10", "100", "00 ", "0 😊"},
			4: {"user", "ser ", "er 1", "r 10", " 100", "100 ", "00 😊"},
		},
	},
	{
		name: "sentence",
		text: "Wait... why? 🤔",
		tokens: map[uint][]string{
			2: {"Wa", "ai", "it", "t.", "..", "..", ". ", " w", "wh", "hy", "y?", "? ", " 🤔"},
			3: {"Wai", "ait", "it.", "t..", "...", ".. ", ". w", " wh", "why", "hy?", "y? ", "? 🤔"},
			4: {"Wait", "ait.", "it..", "t...", "... ", ".. w", ". wh", " why", "why?", "hy? ", "y? 🤔"},
		},
		ctokens: map[uint][]string{
			2: {"Wa", "ai", "it", "t ", " w", "wh", "hy", "y ", " 🤔"},
			3: {"Wai", "ait", "it ", "t w", " wh", "why", "hy ", "y 🤔"},
			4: {"Wait", "ait ", "it w", "t wh", " why", "why ", "hy 🤔"},
		},
	},
	{
		name: "long sentence",
		text: "GitHub (©2024) — awesome! 🚀",
		tokens: map[uint][]string{
			2: {"Gi", "it", "tH", "Hu", "ub", "b ", " (", "(©", "©2", "20", "02", "24", "4)", ") ", " —", "— ", " a", "aw", "we", "es", "so", "om", "me", "e!", "! ", " 🚀"},
			3: {"Git", "itH", "tHu", "Hub", "ub ", "b (", " (©", "(©2", "©20", "202", "024", "24)", "4) ", ") —", " — ", "— a", " aw", "awe", "wes", "eso", "som", "ome", "me!", "e! ", "! 🚀"},
			4: {"GitH", "itHu", "tHub", "Hub ", "ub (", "b (©", " (©2", "(©20", "©202", "2024", "024)", "24) ", "4) —", ") — ", " — a", "— aw", " awe", "awes", "weso", "esom", "some", "ome!", "me! ", "e! 🚀"},
		},
		ctokens: map[uint][]string{
			2: {"Gi", "it", "tH", "Hu", "ub", "b ", " ©", "©2", "20", "02", "24", "4 ", "  ", " a", "aw", "we", "es", "so", "om", "me", "e ", " 🚀"},
			3: {"Git", "itH", "tHu", "Hub", "ub ", "b ©", " ©2", "©20", "202", "024", "24 ", "4  ", "  a", " aw", "awe", "wes", "eso", "som", "ome", "me ", "e 🚀"},
			4: {"GitH", "itHu", "tHub", "Hub ", "ub ©", "b ©2", " ©20", "©202", "2024", "024 ", "24  ", "4  a", "  aw", " awe", "awes", "weso", "esom", "some", "ome ", "me 🚀"},
		},
	},
}

func TestChar(t *testing.T) {
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
	for i := 0; i < len(cstages); i++ {
		st := &cstages[i]
		t.Run(st.name, func(t *testing.T) {
			t.Run("origin", func(t *testing.T) {
				for k, list := range st.tokens {
					t.Run(strconv.Itoa(int(k)), func(t *testing.T) {
						sh := NewChar[string](uint64(k), "")
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
						sh := NewChar[string](uint64(k), CleanSetAll)
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

func BenchmarkChar(b *testing.B) {
	for i := 0; i < len(cstages); i++ {
		st := &cstages[i]
		b.Run(st.name, func(b *testing.B) {
			b.Run("origin", func(b *testing.B) {
				for k, _ := range st.tokens {
					b.Run(strconv.Itoa(int(k)), func(b *testing.B) {
						sh := NewChar[string](uint64(k), "")
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
						sh := NewChar[string](uint64(k), CleanSetAll)
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
