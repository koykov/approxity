package pbtk

import (
	"fmt"
	"testing"
)

type (
	thashable[T Hashable]      []T
	thashableStage[T Hashable] struct {
		vals thashable[T]
		res  thashable[T]
	}
)

var (
	thiStages = []thashableStage[int]{
		{
			vals: []int{1, 2, 3, 4, 5},
			res:  []int{1, 2, 3, 4, 5},
		},
		{
			vals: []int{1, 1, 1, 1, 1},
			res:  []int{1},
		},
	}
	thsStages = []thashableStage[string]{
		{
			vals: []string{"a", "b", "c", "d", "e"},
			res:  []string{"a", "b", "c", "d", "e"},
		},
		{
			vals: []string{"a", "a", "a", "a", "a"},
			res:  []string{"a"},
		},
	}
)

func TestDeduplicate(t *testing.T) {
	for _, stage := range thiStages {
		t.Run("", func(t *testing.T) {
			res := Deduplicate(stage.vals)
			if fmt.Sprintf("%v", res) != fmt.Sprintf("%v", stage.res) {
				t.Errorf("Deduplicate(%v) = %v, want %v", stage.vals, res, stage.res)
			}
		})
	}
	for _, stage := range thsStages {
		t.Run("", func(t *testing.T) {
			res := Deduplicate(stage.vals)
			if fmt.Sprintf("%v", res) != fmt.Sprintf("%v", stage.res) {
				t.Errorf("Deduplicate(%v) = %v, want %v", stage.vals, res, stage.res)
			}
		})
	}
}

func BenchmarkDeduplicate(b *testing.B) {
	for _, stage := range thiStages {
		b.Run("", func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				Deduplicate(stage.vals)
			}
		})
	}
	for _, stage := range thsStages {
		b.Run("", func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				Deduplicate(stage.vals)
			}
		})
	}
}
