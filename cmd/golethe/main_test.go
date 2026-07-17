package main

import (
	"testing"

	"github.com/fmportela/golethe/internal/engine"
)

func TestSpaceReleasesWord(t *testing.T) {
	model := &engine.Model{}
	for _, r := range "lethe" {
		handleRune(model, r)
	}
	handleRune(model, ' ')

	if got := model.ActiveText(); got != "" {
		t.Fatalf("active word = %q, want empty", got)
	}
	if got := model.DisplayWords(); len(got) != 1 || got[0] != "lethe" {
		t.Fatalf("display words = %#v, want lethe", got)
	}
}

func TestClearKeysClearTheTrail(t *testing.T) {
	for _, key := range []rune{'\n', 0x0c} {
		t.Run(string(key), func(t *testing.T) {
			model := &engine.Model{}
			for _, r := range "lethe" {
				handleRune(model, r)
			}
			handleRune(model, ' ')
			for _, r := range "go" {
				handleRune(model, r)
			}
			handleRune(model, key)

			if got := model.ActiveText(); got != "" {
				t.Fatalf("active word = %q, want empty", got)
			}
			if got := model.DisplayWords(); len(got) != 0 {
				t.Fatalf("display words = %#v, want none", got)
			}
		})
	}
}

func TestBackspaceRemovesTheLastRune(t *testing.T) {
	model := &engine.Model{}
	for _, r := range "lethe" {
		handleRune(model, r)
	}
	handleRune(model, 0x7f)

	if got := model.ActiveText(); got != "leth" {
		t.Fatalf("active word = %q, want leth", got)
	}
}
