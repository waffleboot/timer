package main

import "time"

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
