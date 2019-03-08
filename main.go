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

type exitRequestError struct{}

type item struct {
	time time.Duration
	name string
}

type term struct {
	time   time.Time
	items  []item
	prompt string
	line   *liner.State
}

func (exitRequestError) Error() string {
	return ""
}

func (t *term) call(cmdstr string) error {
	if strings.HasPrefix(cmdstr, "-") {
		return t.subtime(cmdstr[1:])
	}
	s := strings.Fields(cmdstr)
	switch s[0] {
	case "q", "e":
		return exitRequestError{}
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
	} else if d != 0 {
		t.additem(d, desc)
	}
	return nil
}

func (t *term) additem(d time.Duration, s string) {
	t.items = append(t.items, item{d, s})
	t.cmdshow()
}

func parseduration(s string) (time.Duration, string, error) {
	var x, y int
	a := strings.Fields(s)
	b := strings.SplitN(a[0], ":", 2)
	desc := ""
	if len(a) > 1 {
		desc = a[1]
	}
	x, err := strconv.Atoi(b[0])
	if err != nil {
		return 0, "", err
	}
	if len(b) > 1 {
		if len(b[1]) > 0 {
			y, err = strconv.Atoi(b[1])
			if err != nil {
				return 0, "", nil
			}
		}
		return time.Duration(x)*time.Hour + time.Duration(y)*time.Minute, desc, nil
	}
	return time.Duration(x) * time.Minute, desc, nil
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
	var sum time.Duration
	if len(t.items) > 0 {
		for i, s := range t.items {
			fmt.Printf("%v)\t%v %v\n", i+1, formatduration(s.time), s.name)
			sum += s.time
		}
		fmt.Println("--------------")
		if !t.time.IsZero() {
			fmt.Printf("%v-->%v (%v)\n", formattime(t.time), formattime(t.time.Add(-sum)), formatduration(sum))
		} else {
			fmt.Printf("total\t%v\n", formatduration(sum))
		}
	}
	return nil
}

type customduration struct {
	name string
	dur  time.Duration
	desc string
}

var customdurations []customduration

func init() {
	customdurations = []customduration{
		{"отрадное", time.Duration(30) * time.Minute, "до станции отрадная от дома"},
		{"кунцево", time.Duration(10) * time.Minute, "от станции кунцево до метро"},
		{"молодежная", time.Duration(10) * time.Minute, "до молодежной и обратно на кунцевскую"},
	}
}

func completeCustomDuration(t string) []string {
	if len(t) > 0 {
		s := strings.Fields(t)[0]
		for _, custom := range customdurations {
			if strings.HasPrefix(custom.name, s) {
				return []string{"-" + custom.name}
			}
		}
	}
	return nil
}

func showCustomDurations() {
	fmt.Println("custom durations:")
	for _, d := range customdurations {
		fmt.Printf("-%-20s\t%s\n", d.name, d.desc)
	}
}

func selectCustomDuration(s string) *customduration {
	for i := range customdurations {
		if s == customdurations[i].name {
			return &customdurations[i]
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
		prompt: "$ ",
		line:   liner.NewLiner(),
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
			usage()
		} else if err = t.call(cmdstr); err != nil {
			if _, ok := err.(exitRequestError); ok {
				return
			}
			fmt.Fprintf(os.Stderr, "command failed: %s\n", err)
		}
	}
}

func usage() {
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
	fmt.Println("$ press Enter to show usage")
	newTerm().Run()
}

func formattime(t time.Time) string {
	return t.Format("15:04")
}

func formatduration(d time.Duration) string {
	h := d / time.Hour
	m := (d - h*time.Hour) / time.Minute
	if h == 0 {
		return fmt.Sprintf("% 2d:%02d", h, m)
	} else {
		return fmt.Sprintf("%02d:%02d", h, m)
	}
}
