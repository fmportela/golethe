# golethe

[![CI](https://github.com/fmportela/golethe/actions/workflows/ci.yml/badge.svg)](https://github.com/fmportela/golethe/actions/workflows/ci.yml)
[![Latest release](https://img.shields.io/github/v/release/fmportela/golethe)](https://github.com/fmportela/golethe/releases)
[![License](https://img.shields.io/github/license/fmportela/golethe)](LICENSE)
[![Go version](https://img.shields.io/badge/go-1.24.0-00ADD8?logo=go)](go.mod)

*lethe it go*

In Greek, *lethe* means forgetfulness or oblivion. It alludes to the
River Lethe of Greek mythology, whose waters made the dead forget their past
lives. `golethe` borrows that idea: words are written briefly, then released
and allowed to disappear.

This is a hobby project built to learn Go fundamentals, including package
structure, goroutines, channels, terminal input, ANSI escape sequences, and
testing. Parts of the code were developed with AI assistance as part of that
learning process.

## Run

```bash
make run
```

Choose a different trail length with `go run ./cmd/golethe -words 8`.

## Install A Release

GitHub Releases provide prebuilt archives for Linux and macOS on both `amd64`
and `arm64`. Download the archive matching your operating system and CPU,
extract it, then run:

```bash
./golethe
```

No Go installation or Makefile is required. New releases are created when a
Git tag beginning with `v` is pushed, for example `v1.0.0`.

## Controls

| Key | Action |
| --- | --- |
| Printable character | Add it to the active word |
| Backspace | Remove the final character |
| Space | Release the active word into the trail |
| Enter / Ctrl-L | Clear the active word and trail |
| Esc / Ctrl-C | Exit cleanly |

By default, the 10 most recently released words are retained. The next word
removes the oldest immediately. Nothing is written to disk.
