package main

import (
	"strings"

	"github.com/peterh/liner"
)

type term struct {
	liner  *liner.State
	prompt string
}

func newTerm() *term {
	t := &term{
		liner:  liner.NewLiner(),
		prompt: "$ ",
	}
	t.liner.SetCompleter(func(line string) []string {
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
	l, err := t.liner.Prompt(t.prompt)
	if err != nil {
		return "", err
	}
	l = strings.TrimSuffix(l, "\n")
	if l != "" {
		t.liner.AppendHistory(l)
	}
	return l, nil
}

func (t *term) close() {
	t.liner.Close()
}
