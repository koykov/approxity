package shingle

import "testing"

func TestHbtrie(t *testing.T) {
	var trie hbtrie
	for _, r := range CleanSetAll {
		trie.set(r)
	}
	for _, r := range CleanSetAll {
		if !trie.contains(r) {
			t.Errorf("trie should contain %c", r)
		}
	}
	for _, r := range "abcdefghijklmnopqrstuvwxyz" {
		if trie.contains(r) {
			t.Errorf("trie should not contain %c", r)
		}
	}
	for _, r := range "абырвалг" {
		if trie.contains(r) {
			t.Errorf("trie should not contain %c", r)
		}
	}
	trie.reset()
	for _, r := range CleanSetAll {
		if trie.contains(r) {
			t.Errorf("trie should not contain %c", r)
		}
	}
}

func BenchmarkHbtrie(b *testing.B) {
	var trie hbtrie
	for _, r := range CleanSetAll {
		trie.set(r)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		trie.contains('.')
	}
}
