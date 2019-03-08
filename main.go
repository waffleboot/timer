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

func (exitRequestError) Error() string {
	return ""
}

func (t *term) call(cmdstr string) error {
	if strings.HasPrefix(cmdstr, "-") {
		return t.subtime(cmdstr[1:])
	}
	s := strings.Fields(cmdstr)
	switch s[0] {
	case "e":
		return exitRequestError{}
	case "s":
		return t.cmdshow()
	case "d":
		return t.cmddel(s[1:])
	case "t":
		return t.cmdtime(s[1:])
	}
	return errors.New(cmdstr)
}

type path struct {
	time time.Duration
	name string
}

type term struct {
	starttime   time.Time
	subtraction []path
	prompt      string
	line        *liner.State
}

func formattime(t time.Time) string {
	return t.Format("15:04")
}

func formatduration(d time.Duration) string {
	h := d / time.Hour
	m := (d - h*time.Hour) / time.Minute
	return fmt.Sprintf("%02d:%02d", h, m)
}

func (t *term) cmddel(s []string) error {
	if len(s) == 0 {
		t.subtraction = t.subtraction[0:0]
		return nil
	}
	p, err := strconv.Atoi(s[0])
	if err != nil || p < 1 || len(t.subtraction) < p {
		return nil
	}
	t.subtraction = append(t.subtraction[:p-1], t.subtraction[p:]...)
	t.cmdshow()
	return nil
}

func (t *term) cmdshow() error {
	var i int
	var sum time.Duration
	if len(t.subtraction) > 0 {
		for _, s := range t.subtraction {
			i++
			fmt.Printf("%v)\t%v %v\n", i, formatduration(s.time), s.name)
			sum += s.time
		}
		fmt.Println("--------------")
		if !t.starttime.IsZero() {
			fmt.Printf("%v-->%v (%v)\n", formattime(t.starttime), formattime(t.starttime.Add(-sum)), formatduration(sum))
		} else {
			fmt.Printf("total\t%v\n", formatduration(sum))
		}
	}
	return nil
}

func parseduration(s string) (time.Duration, string, error) {
	var x, y int
	a := strings.Fields(s)
	b := strings.SplitN(a[0], ":", 2)
	n := ""
	if len(a) > 1 {
		n = a[1]
	}
	x, err := strconv.Atoi(b[0])
	if err != nil {
		return 0, "", err
	}
	if len(b) > 1 {
		y, err = strconv.Atoi(b[1])
		if err != nil {
			return 0, "", nil
		}
		return time.Duration(x)*time.Hour + time.Duration(y)*time.Minute, n, nil
	}
	return time.Duration(x) * time.Minute, n, nil
}

func (t *term) addtime(d time.Duration, s string) {
	t.subtraction = append(t.subtraction, path{d, s})
	t.cmdshow()
}

var customdurations []struct {
	name string
	dur  time.Duration
	desc string
}

func init() {
	customdurations = []struct {
		name string
		dur  time.Duration
		desc string
	}{
		{"отрадное", time.Duration(30) * time.Minute, "до станции отрадная от дома"},
		{"кунцево", time.Duration(10) * time.Minute, "от станции кунцево до метро"},
		{"молодежная", time.Duration(10) * time.Minute, "до молодежной и обратно на кунцевскую"},
	}
}

func (t *term) subtime(s string) error {
	for _, custom := range customdurations {
		if s == custom.name {
			t.addtime(custom.dur, custom.desc)
			return nil
		}
	}
	if s == "" {
		fmt.Fprintf(os.Stderr, "custom durations:\n")
		for _, d := range customdurations {
			fmt.Fprintf(os.Stderr, "-%-20s\t%s\n", d.name, d.desc)
		}
		return nil
	}
	d, n, err := parseduration(s)
	if err != nil {
		return err
	} else if d != 0 {
		t.addtime(d, n)
	}
	return nil
}

func parsetime2(s []string) (time.Time, error) {
	if len(s) == 1 {
		return time.Parse("15", s[0])
	}
	return time.Parse("15:04", s[0]+":"+s[1])
}

func parsetime(s []string) (time.Time, error) {
	n, err := time.Parse("15:04", s[0])
	if err != nil {
		return parsetime2(s)
	}
	return n, nil
}

func (t *term) settime(s []string) error {
	n, err := parsetime(s)
	if err != nil {
		return err
	}
	t.starttime = n
	return nil
}

func (t *term) cmdtime(s []string) error {
	if len(s) > 0 {
		t.settime(s)
		t.cmdshow()
	}
	return nil
}

func newTerm() *term {
	return &term{
		prompt: "$ ",
		line:   liner.NewLiner(),
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, `		commands:
[e]xit
[s]how		show full timetable
[d]elete n	delete n-th item
[t]ime		set timepoint
-		show custom durations
-mm		add time interval
-hh:mm		add time interval
-name		add custom interval
`)
}

func (t *term) Run() {
	defer t.close()
	t.line.SetCompleter(func(line string) []string {
		if strings.HasPrefix(line, "-") {
			t := strings.TrimSpace(line[1:])
			if len(t) > 0 {
				s := strings.Fields(t)[0]
				for _, custom := range customdurations {
					if strings.HasPrefix(custom.name, s) {
						return []string{"-" + custom.name}
					}
				}
			}
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
	fmt.Fprintf(os.Stderr, "$ press Enter to show usage\n")
	newTerm().Run()
}
