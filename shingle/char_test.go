package shingle

import (
	"strconv"
	"testing"
)

type stage struct {
	name    string
	text    string
	ngrams  map[uint][]string
	cngrams map[uint][]string
}

var stages = []stage{
	{
		name: "single word",
		text: "Hello!",
		ngrams: map[uint][]string{
			2: {"He", "el", "ll", "lo", "o!"},
			3: {"Hel", "ell", "llo", "lo!"},
			4: {"Hell", "ello", "llo!"},
		},
		cngrams: map[uint][]string{
			2: {"He", "el", "ll", "lo"},
			3: {"Hel", "ell", "llo"},
			4: {"Hell", "ello"},
		},
	},
	{
		name: "phrase",
		text: "@user: $100 😊",
		ngrams: map[uint][]string{
			2: {"@u", "us", "se", "er", "r:", ": ", " $", "$1", "10", "00", "0 ", " 😊"},
			3: {"@us", "use", "ser", "er:", "r: ", ": $", " $1", "$10", "100", "00 ", "0 😊"},
			4: {"@use", "user", "er:", "r: ", ": $", " $1", "$10", "$100", "100 ", "00 😊"},
		},
		cngrams: map[uint][]string{
			2: {"us", "se", "er", "r ", " 1", "10", "00", "0 ", " 😊"},
			3: {"use", "ser", "er ", "r 1", " 10", "100", "00 ", "0 😊"},
			4: {"user", "ser ", "er 1", "r 10", " 100", "100 ", "00 😊"},
		},
	},
	{
		name: "sentence",
		text: "Wait... why? 🤔",
		ngrams: map[uint][]string{
			2: {"Wa", "ai", "it", "t.", ".", ".", ". ", " w", "wh", "hy", "y?", "? ", " 🤔"},
			3: {"Wai", "ait", "it.", "t..", "...", ".. ", ". w", " wh", "why", "hy?", "y? ", "? 🤔"},
			4: {"Wait", "ait.", "it..", "t...", "... ", ".. w", ". wh", " why", "why?", "hy? ", "y? 🤔"},
		},
		cngrams: map[uint][]string{
			2: {"Wa", "ai", "it", "t ", " w", "wh", "hy", "y ", " 🤔"},
			3: {"Wai", "ait", "it ", "t w", " wh", "why", "hy ", "y 🤔"},
			4: {"Wait", "ait ", "it w", "t wh", " why", "why ", "hy 🤔"},
		},
	},
	{
		name: "long sentence",
		text: "GitHub (©2024) — awesome! 🚀",
		ngrams: map[uint][]string{
			2: {"Gi", "it", "tH", "Hu", "ub", "b ", " (", "(©", "©2", "20", "02", "24", "4)", "),", " —", "— ", " a", "aw", "we", "es", "so", "om", "me", "e!", "! ", " 🚀"},
			3: {"Git", "itH", "tHu", "Hub", "ub ", "b (", " (©", "(©2", "©20", "202", "024", "24)", "4),", ") —", "— a", " aw", "awe", "wes", "eso", "som", "ome", "me!", "e! ", "! 🚀"},
			4: {"GitH", "itHu", "tHub", "Hub ", "ub (", "b (©", " (©2", "(©20", "©202", "2024", "024)", "24),", "4) —", ") — ", "— aw", " awe", "awes", "weso", "esom", "some", "ome!", "me! ", "e! 🚀"},
		},
		cngrams: map[uint][]string{
			2: {"Gi", "it", "tH", "Hu", "ub", "b ", "©2", "20", "02", "24", "4 ", " a", "aw", "we", "es", "so", "om", "me", "e ", " 🚀"},
			3: {"Git", "itH", "tHu", "Hub", "ub ", "b ©", "©20", "202", "024", "24 ", " aw", "awe", "wes", "eso", "som", "ome", "me ", " 🚀"},
			4: {"GitH", "itHu", "tHub", "Hub ", "ub ©", "b ©2", "©202", "2024", "024 ", " awe", "awes", "weso", "esom", "some", "ome ", " 🚀"},
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
	t.Run("origin", func(t *testing.T) {
		for i := 0; i < len(stages); i++ {
			st := &stages[i]
			t.Run(st.name, func(t *testing.T) {
				t.Run("origin", func(t *testing.T) {
					for k, list := range st.ngrams {
						t.Run(strconv.Itoa(int(k)), func(t *testing.T) {
							sh := NewChar[string](k, false)
							r := sh.Shingle(st.text)
							if !sheq(r, list) {
								t.Errorf("expected %v, got %v", list, r)
							}
						})
					}
				})
				t.Run("clean", func(t *testing.T) {
					for k, list := range st.cngrams {
						t.Run(strconv.Itoa(int(k)), func(t *testing.T) {
							sh := NewChar[string](k, true)
							r := sh.Shingle(st.text)
							if !sheq(r, list) {
								t.Errorf("expected %v, got %v", list, r)
							}
						})
					}
				})
			})
		}
	})
}
