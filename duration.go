package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type duration struct {
	hours   int
	minutes int
}

func (d duration) isValid() bool {
	return d.hours != 0 || d.minutes != 0
}

func (d duration) add(o duration) duration {
	t := time.Duration(d.hours)*time.Hour + time.Duration(d.minutes)*time.Minute
	t += time.Duration(o.hours)*time.Hour + time.Duration(o.minutes)*time.Minute
	h := t / time.Hour
	m := (t - h*time.Hour) / time.Minute
	return duration{int(h), int(m)}
}

func sum(t time.Time, d duration) time.Time {
	return t.Add(-time.Duration(d.hours)*time.Hour - time.Duration(d.minutes)*time.Minute)
}

func (d duration) Format(f fmt.State, c rune) {
	_, ok := f.Width()
	if ok {
		fmt.Fprintf(f, "%2d:%02d", d.hours, d.minutes)
	} else {
		fmt.Fprintf(f, "%d:%02d", d.hours, d.minutes)
	}
}

func parseduration(s string) (duration, string, error) {
	var x, y int
	a := strings.Fields(s)
	b := strings.SplitN(a[0], ":", 2)
	desc := ""
	if len(a) > 1 {
		desc = strings.Join(a[1:], " ")
	}
	x, err := strconv.Atoi(b[0])
	if err != nil {
		return duration{0, 0}, "", err
	}
	if len(b) > 1 {
		if len(b[1]) > 0 {
			y, err = strconv.Atoi(b[1])
			if err != nil {
				return duration{0, 0}, "", nil
			}
		}
		return duration{x, y}, desc, nil
	}
	return duration{0, x}, desc, nil
}
