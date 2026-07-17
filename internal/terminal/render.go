// Package terminal handles raw mode and ANSI rendering for a Unix terminal.
package terminal

import (
	"fmt"
	"io"
	"strings"
	"unicode/utf8"

	"github.com/fmportela/golethe/internal/engine"
)

const (
	brightGray = 255
	darkGray   = 232
)

type cell struct {
	row    int
	column int
	width  int
}

// Renderer draws the recent word trail and clears only its previous position.
type Renderer struct {
	out      io.Writer
	previous []cell
}

// NewRenderer creates a renderer that writes ANSI frames to out.
func NewRenderer(out io.Writer) *Renderer {
	return &Renderer{out: out}
}

// Draw renders the model at the terminal's current size.
func (r *Renderer) Draw(model *engine.Model, rows, columns int) {
	if rows < 1 || columns < 1 {
		return
	}

	var frame strings.Builder
	frame.WriteString(home)
	for _, previous := range r.previous {
		frame.WriteString(cursor(previous.row, previous.column))
		frame.WriteString(strings.Repeat(" ", previous.width))
	}

	// Keep the writing trail in the upper portion of the otherwise empty screen.
	horizon := max(1, rows/3)
	words := model.DisplayWords()
	active := model.ActiveText()
	for index, word := range words {
		// Reserve the final column for the terminal cursor.
		words[index] = truncateToWidth(word, max(1, columns-1))
	}
	if len(words) == 0 && active == "" {
		frame.WriteString(reset)
		frame.WriteString(cursor(horizon, max(1, columns/2)))
		fmt.Fprint(r.out, frame.String())
		r.previous = nil
		return
	}

	// Keep the cursor still at the center. The active word grows to its left.
	cursorColumn := max(1, columns/2)
	focus := ""
	focusWidth := 0
	focusColumn := cursorColumn
	history := words
	colorTotal := len(words)
	if active != "" {
		focus = truncateFromLeft(words[len(words)-1], max(1, cursorColumn-1))
		focusWidth = utf8.RuneCountInString(focus)
		focusColumn = cursorColumn - focusWidth
		history = words[:len(words)-1]
	}
	historyStart := 0
	availableHistoryWidth := max(0, focusColumn-2)
	for historyStart < len(history) && lineWidth(history[historyStart:]) > availableHistoryWidth {
		historyStart++
	}
	history = history[historyStart:]
	historyWidth := lineWidth(history)
	historyColumn := max(1, focusColumn-historyWidth-1)
	column := historyColumn
	for index, word := range history {
		frame.WriteString(cursor(horizon, column))
		frame.WriteString(color(wordColor(historyStart+index, colorTotal)))
		frame.WriteString(word)
		column += utf8.RuneCountInString(word) + 1
	}
	if focus != "" {
		frame.WriteString(cursor(horizon, focusColumn))
		frame.WriteString(color(wordColor(len(words)-1, colorTotal)))
		frame.WriteString(focus)
	}

	frame.WriteString(reset)
	frame.WriteString(cursor(horizon, cursorColumn))
	fmt.Fprint(r.out, frame.String())
	r.previous = []cell{{
		row:    horizon,
		column: max(1, historyColumn),
		width:  cursorColumn - max(1, historyColumn),
	}}
}

// lineWidth returns the number of rune columns occupied by words and spaces.
func lineWidth(words []string) int {
	if len(words) == 0 {
		return 0
	}

	width := len(words) - 1
	for _, word := range words {
		width += utf8.RuneCountInString(word)
	}
	return width
}

// truncateToWidth shortens text to width runes, adding an ellipsis when possible.
func truncateToWidth(text string, width int) string {
	runes := []rune(text)
	if len(runes) <= width {
		return text
	}
	if width <= 3 {
		return string(runes[:width])
	}
	return string(runes[:width-3]) + "..."
}

// truncateFromLeft keeps the end of text visible, adding an ellipsis when possible.
func truncateFromLeft(text string, width int) string {
	runes := []rune(text)
	if len(runes) <= width {
		return text
	}
	if width <= 3 {
		return string(runes[len(runes)-width:])
	}
	return "..." + string(runes[len(runes)-(width-3):])
}

// wordColor maps a word's oldest-to-newest position onto the grayscale ramp.
func wordColor(index, total int) int {
	if total <= 1 {
		return brightGray
	}

	return brightGray - (total-1-index)*(brightGray-darkGray)/(total-1)
}
