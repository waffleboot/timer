package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/peterh/liner"
)

var exitRequestError error = errors.New("")

type item struct {
	time duration
	name string
}

type term struct {
	line   *liner.State
	prompt string

	time   time.Time
	items  []item
}

func (t *term) parseCommandText(cmdstr string) error {
	if strings.HasPrefix(cmdstr, "-") {
		return t.subtime(cmdstr[1:])
	}
	s := strings.Fields(cmdstr)
	switch s[0] {
	case "q", "e":
		return exitRequestError
	case "s":
		return t.cmdshow()
	case "d":
		return t.cmddel(s[1:])
	case "t":
		return t.cmdtime(s[1:])
	default:
		return errors.New(cmdstr)
	}
}

func (t *term) subtime(s string) error {
	if s == "" {
		showCustomDurations()
		return nil
	}
	custom := selectCustomDuration(s)
	if custom != nil {
		t.additem(custom.dur, custom.desc)
		return nil
	}
	d, desc, err := parseduration(s)
	if err != nil {
		return err
	} else if d.isValid() {
		t.additem(d, desc)
	}
	return nil
}

func (t *term) additem(d duration, s string) {
	t.items = append(t.items, item{d, s})
	t.cmdshow()
}

func (t *term) cmddel(s []string) error {
	if len(s) == 0 {
		t.items = t.items[:0]
		fmt.Println("all items deleted")
		return nil
	}
	p, err := strconv.Atoi(s[0])
	if err != nil || p < 1 || len(t.items) < p {
		return nil
	}
	t.items = append(t.items[:p-1], t.items[p:]...)
	return t.cmdshow()
}

func (t *term) cmdshow() error {
	var total duration
	if len(t.items) > 0 {
		for i, s := range t.items {
			total = total.add(s.time)
			fmt.Printf("%v)\t%2v\t%2v %v\n", i+1, s.time, formattime(sum(t.time, total)), s.name)
		}
		fmt.Println("---------------------")
		if !t.time.IsZero() {
			fmt.Printf("%v-->%v  (%v)\n", formattime(t.time), formattime(sum(t.time, total)), total)
		} else {
			fmt.Printf("total\t%v\n", total)
		}
	}
	return nil
}

func (t *term) cmdtime(s []string) error {
	if len(s) > 0 {
		t.settime(s)
		t.cmdshow()
	}
	return nil
}

func (t *term) settime(s []string) error {
	time, err := parsehhmm(s)
	if err != nil {
		return err
	}
	t.time = time
	return nil
}

func parsehhmm(s []string) (time.Time, error) {
	time, err := time.Parse("15:04", s[0])
	if err != nil {
		return parsehh(s)
	}
	return time, err
}

func parsehh(s []string) (time.Time, error) {
	if len(s) == 1 {
		t, err := time.Parse("15", s[0])
		if err != nil {
			return time.Parse("15:", s[0])
		} else {
			return t, err
		}
	}
	return time.Parse("15:04", s[0]+":"+s[1])
}

//-----------------------------------------------------------------------------

func newTerm() *term {
	return &term{
		line:   liner.NewLiner(),
		prompt: "$ ",
	}
}

func (t *term) Run() {
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
		} else if err = t.parseCommandText(cmdstr); err != nil {
			if err == exitRequestError {
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

func main() {
	printUsage()
	newTerm().Run()
}

func formattime(t time.Time) string {
	return t.Format("15:04")
}

func printUsage() {
	fmt.Printf(`		commands:
[q]uit
[e]xit
[s]how		show full timetable
[d]elete	delete all items
[d]elete n	delete n-th item
[t]ime hh[:mm]	set timepoint
-		show custom durations
-mm		add time interval
-hh[:mm]	add time interval
-name		add custom interval

`)
}
