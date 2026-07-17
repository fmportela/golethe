package engine

import (
	"fmt"
	"testing"
)

func TestReleaseKeepsOnlyNewestWords(t *testing.T) {
	model := &Model{}
	for number := range DefaultMaxWords + 1 {
		text := fmt.Sprintf("word%d", number)
		for _, r := range text {
			model.AppendRune(r)
		}
		model.Release()
	}

	words := model.DisplayWords()
	if len(words) != DefaultMaxWords {
		t.Fatalf("got %d words, want %d", len(words), DefaultMaxWords)
	}
	if words[0] != "word1" || words[len(words)-1] != fmt.Sprintf("word%d", DefaultMaxWords) {
		t.Fatalf("kept %q through %q, want word1 through word%d", words[0], words[len(words)-1], DefaultMaxWords)
	}
}

func TestNewModelUsesConfiguredWordLimit(t *testing.T) {
	model := NewModel(2)
	for _, text := range []string{"one", "two", "three"} {
		for _, r := range text {
			model.AppendRune(r)
		}
		model.Release()
	}

	words := model.DisplayWords()
	if len(words) != 2 || words[0] != "two" || words[1] != "three" {
		t.Fatalf("display words = %#v, want two and three", words)
	}
}

func TestDisplayWordsIncludesTheActiveWord(t *testing.T) {
	model := &Model{}
	for _, r := range "lethe" {
		model.AppendRune(r)
	}
	model.Release()
	for _, r := range "go" {
		model.AppendRune(r)
	}

	words := model.DisplayWords()
	if len(words) != 2 || words[0] != "lethe" || words[1] != "go" {
		t.Fatalf("display words = %#v, want lethe and go", words)
	}
}
