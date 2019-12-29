package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/peterh/liner"
)

type term struct {
	line   *liner.State
	prompt string
}

func newTerm() *term {
	return &term{
		line:   liner.NewLiner(),
		prompt: "$ ",
	}
}

func (t *term) Run(tg *timing) {
	defer t.close()
	t.line.SetCompleter(func(line string) []string {
		if strings.HasPrefix(line, "-") {
			return completeCustomDuration(strings.TrimSpace(line[1:]))
		}
		return nil
	})
	for {
		cmdstr, err := t.promptForInput()
		if err != nil {
			return
		} else if cmdstr == "" {
			printUsage()
		} else if err = tg.parseCommandText(cmdstr); err != nil {
			if err == errExitRequest {
				break
			}
			fmt.Fprintf(os.Stderr, "command failed: %s\n", err)
		}
	}
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
