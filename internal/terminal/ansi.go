// Package terminal handles raw mode and ANSI rendering for a Unix terminal.
package terminal

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const (
	esc        = "\x1b"        // Escape character that begins an ANSI sequence.
	clear      = esc + "[2J"   // Clear the terminal's visible screen.
	home       = esc + "[H"    // Move the cursor to row 1, column 1.
	hideCursor = esc + "[?25l" // Hide the cursor.
	showCursor = esc + "[?25h" // Show the cursor.
	reset      = esc + "[0m"   // Reset colors and other text attributes.
)

// Terminal owns raw-mode setup and the ANSI sequences used by the renderer.
type Terminal struct {
	out       *os.File
	sttyState string
}

// New saves the current terminal settings, then enables raw input without echo.
// Call Restore when the terminal is no longer in use.
func New(out *os.File) (*Terminal, error) {
	state, err := runStty("-g")
	if err != nil {
		return nil, fmt.Errorf("read terminal settings: %w", err)
	}
	if _, err := runStty("raw", "-echo"); err != nil {
		return nil, fmt.Errorf("enable raw mode: %w", err)
	}
	return &Terminal{out: out, sttyState: state}, nil
}

// Restore reapplies the terminal settings that New saved.
func (t *Terminal) Restore() error {
	if t.sttyState == "" {
		return nil
	}
	_, err := runStty(t.sttyState) // restores original state
	t.sttyState = ""
	return err
}

// Start clears the terminal and prepares it for interactive rendering.
func (t *Terminal) Start() {
	fmt.Fprint(t.out, clear, home, showCursor)
}

// Finish resets terminal formatting, shows the cursor, and clears the screen.
func (t *Terminal) Finish() {
	fmt.Fprint(t.out, reset, showCursor, clear, home)
}

// Size returns the current terminal dimensions as rows and columns.
func (t *Terminal) Size() (rows, columns int, err error) {
	output, err := runStty("size") // returns "rows columns"
	if err != nil {
		return 0, 0, fmt.Errorf("read terminal size: %w", err)
	}
	fields := strings.Fields(output) // converts "rows columns" into []string{"rows", "columns"}
	if len(fields) != 2 {
		return 0, 0, fmt.Errorf("unexpected terminal size %q", output)
	}
	rows, err = strconv.Atoi(fields[0]) // ASCII to integer conversion e.g. "24" -> 24
	if err != nil {
		return 0, 0, fmt.Errorf("parse terminal rows: %w", err)
	}
	columns, err = strconv.Atoi(fields[1])
	if err != nil {
		return 0, 0, fmt.Errorf("parse terminal columns: %w", err)
	}
	return rows, columns, nil
}

// runStty executes the stty command with args attached to the same terminal
// input used by golethe. This lets stty inspect or modify that terminal's
// settings. It returns stty's standard output with surrounding whitespace removed.
func runStty(args ...string) (string, error) {
	command := exec.Command("stty", args...) // prepare the command to run
	command.Stdin = os.Stdin                 // stty reads from standard input
	output, err := command.Output()          // run the command and capture its output
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// cursor returns the ANSI sequence that moves the cursor to one-based coordinates.
func cursor(row, column int) string {
	return fmt.Sprintf("%s[%d;%dH", esc, row, column)
}

// color returns the ANSI sequence that selects a 256-color foreground color.
func color(code int) string {
	return fmt.Sprintf("%s[38;5;%dm", esc, code)
}
