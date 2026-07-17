package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/fmportela/golethe/internal/engine"
	"github.com/fmportela/golethe/internal/terminal"
)

// main parses command-line options and reports any top-level application error.
func main() {
	wordLimit := flag.Int("words", engine.DefaultMaxWords, "number of words to retain")
	flag.Parse()
	if *wordLimit < 1 {
		fmt.Fprintln(os.Stderr, "golethe: -words must be at least 1")
		os.Exit(2)
	}
	if err := run(*wordLimit); err != nil {
		fmt.Fprintln(os.Stderr, "golethe:", err)
		os.Exit(1)
	}
}

// run configures the terminal and processes input, resize, and exit events.
func run(wordLimit int) error {
	term, err := terminal.New(os.Stdout)
	if err != nil {
		return err
	}
	defer term.Restore()
	defer term.Finish()

	term.Start()
	input := readInput(os.Stdin)
	interrupts := make(chan os.Signal, 1)                    // channel to receive interrupt signals (e.g., Ctrl-C)
	signal.Notify(interrupts, os.Interrupt, syscall.SIGTERM) // notify the channel
	defer signal.Stop(interrupts)
	resizes := make(chan os.Signal, 1)       // channel to receive terminal resize signals (SIGWINCH)
	signal.Notify(resizes, syscall.SIGWINCH) // notify the channel
	defer signal.Stop(resizes)

	model := engine.NewModel(wordLimit)
	renderer := terminal.NewRenderer(os.Stdout)
	rows, columns, err := term.Size()
	if err != nil {
		return err
	}
	renderer.Draw(model, rows, columns)

	for {
		select {
		case <-interrupts:
			// if the user presses Ctrl-C or the process receives a termination signal, exit gracefully
			return nil
		case event := <-input:
			// If input has no error, handle the rune and redraw the screen.
			if event.err == nil {
				if handleRune(model, event.rune) {
					return nil
				}
				renderer.Draw(model, rows, columns)
				continue
			}
			if errors.Is(event.err, io.EOF) {
				return nil
			}
			return fmt.Errorf("read input: %w", event.err)
		case <-resizes:
			// if the terminal is resized, get the new size and redraw the screen
			rows, columns, err = term.Size()
			if err != nil {
				return err
			}
			renderer.Draw(model, rows, columns)
		}
	}
}

type inputEvent struct {
	rune rune
	err  error
}

// readInput starts a background goroutine that waits for keyboard input and
// sends decoded runes or one read error through a channel.
func readInput(reader io.Reader) <-chan inputEvent {
	events := make(chan inputEvent)
	go func() {
		input := bufio.NewReader(reader)
		for {
			r, _, err := input.ReadRune() // A rune is Go's representation of a Unicode character
			if err != nil {
				events <- inputEvent{err: err}
				return
			}
			events <- inputEvent{rune: r}
		}
	}()
	return events
}

// handleRune returns true when Ctrl-C requests an exit.
func handleRune(model *engine.Model, r rune) bool {
	switch r {
	case 0x03: // Ctrl-C
		return true
	case ' ':
		model.Release()
	case '\r', '\n', 0x0c: // Enter and Ctrl-L
		model.Clear()
	case 0x7f, 0x08:
		model.Backspace()
	case '\t':
		// Only spaces release words.
	default:
		if r >= 0x20 && r != 0x7f {
			model.AppendRune(r)
		}
	}
	return false
}
