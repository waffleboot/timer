package main

import (
	"strings"

	"github.com/peterh/liner"
)

type term struct {
	line   *liner.State
	prompt string
}

func newTerm() *term {
	t := &term{
		line:   liner.NewLiner(),
		prompt: "$ ",
	}
	t.line.SetCompleter(func(line string) []string {
		if strings.HasPrefix(line, "-") {
			return completeCustomDuration(strings.TrimSpace(line[1:]))
		}
		return nil
	})
	return t
}

func (t *term) read() (string, error) {
	return t.promptForInput()
}

func (t *term) promptForInput() (string, error) {
	l, err := t.line.Prompt(t.prompt)
	if err != nil {
		return "", err
	}
	l = strings.TrimSuffix(l, "\n")
	if l != "" {
		t.line.AppendHistory(l)
	}
	return l, nil
}

func (t *term) close() {
	t.line.Close()
}
