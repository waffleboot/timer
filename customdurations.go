package main

import (
	"fmt"
	"github.com/waffleboot/timer/domain"
	"strings"
)

type customduration struct {
	name string
	dur  domain.Duration
	desc string
}

var customdurations []customduration

func init() {
	customdurations = []customduration{
		{"кунцево", domain.Duration{0, 10}, "до кунцевской"},
		{"отрадное", domain.Duration{0, 30}, "до отрадного"},
		{"электричка", domain.Duration{0, 21}, "на электричке"},
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
