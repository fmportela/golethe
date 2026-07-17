// Package engine owns the terminal-independent state of the void.
package engine

import "strings"

const DefaultMaxWords = 10

// Model contains the word currently being typed and its recent predecessors.
type Model struct {
	active []rune
	words  []string
	limit  int
}

// NewModel creates a word-trail model with the supplied positive word limit.
// Invalid limits use DefaultMaxWords instead.
func NewModel(limit int) *Model {
	if limit < 1 {
		limit = DefaultMaxWords
	}
	return &Model{limit: limit}
}

// ActiveText returns the word currently being typed.
func (m *Model) ActiveText() string {
	return string(m.active)
}

// AppendRune adds one character to the active word.
func (m *Model) AppendRune(r rune) {
	m.active = append(m.active, r)
}

// Backspace removes the final character from the active word, when present.
func (m *Model) Backspace() {
	if len(m.active) > 0 {
		m.active = m.active[:len(m.active)-1]
	}
}

// Release adds the active word to the trail. Empty input has no effect.
func (m *Model) Release() {
	word := strings.TrimSpace(string(m.active))
	m.active = nil
	if word == "" {
		return
	}

	m.words = append(m.words, word)
	if len(m.words) > m.wordLimit() {
		m.words = m.words[len(m.words)-m.wordLimit():]
	}
}

// DisplayWords returns, from oldest to newest, the words visible to the user.
// The active word is included as the newest item while it is being typed.
func (m *Model) DisplayWords() []string {
	words := append([]string(nil), m.words...)
	if active := m.ActiveText(); active != "" {
		words = append(words, active)
	}
	if len(words) > m.wordLimit() {
		words = words[len(words)-m.wordLimit():]
	}
	return words
}

// Clear discards the active word and every word in the visible trail.
func (m *Model) Clear() {
	m.active = nil
	m.words = nil
}

// wordLimit returns a valid limit even when Model was created without NewModel.
func (m *Model) wordLimit() int {
	if m.limit < 1 {
		return DefaultMaxWords
	}
	return m.limit
}
