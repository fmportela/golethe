package terminal

import (
	"bytes"
	"strings"
	"testing"

	"github.com/fmportela/golethe/internal/engine"
)

func TestRendererDrawsTrailOnTheHorizon(t *testing.T) {
	model := engine.NewModel(5)
	for _, text := range []string{"one", "two", "three", "four", "five"} {
		for _, r := range text {
			model.AppendRune(r)
		}
		model.Release()
	}

	var output bytes.Buffer
	NewRenderer(&output).Draw(model, 20, 80)
	for _, expected := range []struct {
		column int
		color  int
		word   string
	}{
		{column: 16, color: 232, word: "one"},
		{column: 20, color: 238, word: "two"},
		{column: 24, color: 244, word: "three"},
		{column: 30, color: 250, word: "four"},
		{column: 35, color: 255, word: "five"},
	} {
		fragment := cursor(6, expected.column) + color(expected.color) + expected.word
		if !strings.Contains(output.String(), fragment) {
			t.Errorf("missing rendered word %q at column %d with color %d", expected.word, expected.column, expected.color)
		}
	}
}

func TestRendererKeepsTheCursorCenteredWhileTheActiveWordGrowsLeft(t *testing.T) {
	for _, text := range []string{"g", "go", "lethe"} {
		t.Run(text, func(t *testing.T) {
			model := &engine.Model{}
			for _, r := range text {
				model.AppendRune(r)
			}

			var output bytes.Buffer
			NewRenderer(&output).Draw(model, 20, 80)

			wordColumn := 40 - len([]rune(text))
			if !strings.Contains(output.String(), cursor(6, wordColumn)+color(255)+text) {
				t.Fatalf("active word was not drawn immediately before the fixed cursor")
			}
			if !strings.HasSuffix(output.String(), reset+cursor(6, 40)) {
				t.Fatal("cursor did not remain at the center")
			}
		})
	}
}

func TestRendererReservesTheNextWordPositionAfterRelease(t *testing.T) {
	model := &engine.Model{}
	for _, r := range "go" {
		model.AppendRune(r)
	}
	model.Release()

	var output bytes.Buffer
	NewRenderer(&output).Draw(model, 20, 80)
	if !strings.Contains(output.String(), cursor(6, 37)+color(255)+"go") {
		t.Fatal("released word did not move left to make room for the next word")
	}
	if !strings.Contains(output.String(), cursor(6, 40)) {
		t.Fatal("empty active-word position was not centered")
	}
}

func TestWordColorStaysInTheGrayscaleRamp(t *testing.T) {
	for index := range 8 {
		color := wordColor(index, 8)
		if color < darkGray || color > brightGray {
			t.Errorf("color %d is outside the grayscale ramp", color)
		}
	}
}

func TestTruncateToWidth(t *testing.T) {
	if got := truncateToWidth("temporary", 6); got != "tem..." {
		t.Fatalf("truncated text = %q, want tem...", got)
	}
}

func TestTruncateFromLeft(t *testing.T) {
	if got := truncateFromLeft("temporary", 6); got != "...ary" {
		t.Fatalf("truncated text = %q, want ...ary", got)
	}
}
