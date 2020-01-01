package main

import (
	"errors"
	"fmt"
	"github.com/waffleboot/timer/domain"
	"os"
	"strconv"
	"strings"
	"time"
)

var errExitRequest error = errors.New("exit")

type input interface {
	read() (string, error)
}

type cli struct {
	service
}

func (t *cli) run(r input) {
	for {
		cmdstr, err := r.read()
		if err != nil {
			return
		} else if cmdstr == "" {
			printUsage()
		} else if err = t.parseCommandText(cmdstr); err != nil {
			if err == errExitRequest {
				break
			}
			fmt.Fprintf(os.Stderr, "command failed: %s\n", err)
		}
	}
}

func (t *cli) parseCommandText(cmdstr string) error {
	if strings.HasPrefix(cmdstr, "-") {
		return t.subtime(cmdstr[1:])
	}
	s := strings.Fields(cmdstr)
	switch s[0] {
	case "q", "e":
		return errExitRequest
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

func (t *cli) subtime(s string) error {
	if s == "" {
		showCustomDurations()
		return nil
	}
	custom := selectCustomDuration(s)
	if custom != nil {
		t.additem(custom.dur, custom.desc)
		return nil
	}
	d, desc, err := domain.ParseDuration(s)
	if err != nil {
		return err
	} else {
		t.additem(d, desc)
	}
	return nil
}

func (t *cli) additem(d domain.Duration, s string) {
	t.Add(d, s)
	t.cmdshow()
}

func (t *cli) cmddel(s []string) error {
	if len(s) == 0 {
		t.Clear()
		fmt.Println("all items deleted")
		return nil
	}
	p, err := strconv.Atoi(s[0])
	if err != nil || p < 1 || t.size() < p {
		return nil
	}
	t.Del(p)
	return t.cmdshow()
}

func (t *cli) cmdshow() error {
	state := t.state()
	if len(state.items) == 0 {
		return nil
	}
	for i, s := range state.items {
		fmt.Printf("%v)\t%2v\t%2v %v\n", i+1, s.Duration, formattime(s.Time), s.desc)
	}
	fmt.Println("---------------------")
	if !state.startTime.IsZero() {
		fmt.Printf("%v-->%v  (%v)\n", formattime(state.startTime), formattime(state.finalTime), state.Duration)
	} else {
		fmt.Printf("total\t%v\n", state.Duration)
	}
	return nil
}

func (t *cli) cmdtime(s []string) error {
	if len(s) > 0 {
		t.settime(s)
		t.cmdshow()
	}
	return nil
}

func (t *cli) settime(s []string) error {
	time, err := parsehhmm(s)
	if err != nil {
		return err
	}
	t.service.settime(time)
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
		}
		return t, err
	}
	return time.Parse("15:04", s[0]+":"+s[1])
}

func formattime(t time.Time) string {
	return t.Format("15:04")
}

//-----------------------------------------------------------------------------

func main() {
	printUsage()
	t := newTerm()
	defer t.close()
	(&cli{}).run(t)
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
