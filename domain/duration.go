package domain

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

//Duration ...
type Duration struct {
	Hours   uint
	Minutes uint
}

func (d Duration) isValid() bool {
	return d.Hours != 0 || d.Minutes != 0
}

//AddDuration ...
func (d Duration) AddDuration(o Duration) Duration {
	t := d.TimeDuration() + o.TimeDuration()
	h := (t / time.Hour)
	m := (t - h*time.Hour) / time.Minute
	return Duration{uint(h), uint(m)}
}

//TimeDuration ...
func (d Duration) TimeDuration() time.Duration {
	return time.Duration(d.Hours)*time.Hour + time.Duration(d.Minutes)*time.Minute
}

////////////////////////////////////////////////////////////////////////////////////

//Format ...
func (d Duration) Format(f fmt.State, c rune) {
	_, ok := f.Width()
	if ok {
		fmt.Fprintf(f, "%2d:%02d", d.Hours, d.Minutes)
	} else {
		fmt.Fprintf(f, "%d:%02d", d.Hours, d.Minutes)
	}
}

//ParseDuration ...
func ParseDuration(s string) (Duration, string, error) {
	var x, y int
	a := strings.Fields(s)
	b := strings.SplitN(a[0], ":", 2)
	desc := ""
	if len(a) > 1 {
		desc = strings.Join(a[1:], " ")
	}
	x, err := strconv.Atoi(b[0])
	if err != nil {
		return Duration{0, 0}, "", err
	}
	if len(b) > 1 {
		if len(b[1]) > 0 {
			y, err = strconv.Atoi(b[1])
			if err != nil {
				return Duration{0, 0}, "", nil
			}
		}
		return Duration{uint(x), uint(y)}, desc, nil
	}
	return Duration{0, uint(x)}, desc, nil
}
