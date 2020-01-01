package main

import (
	"time"

	"github.com/waffleboot/timer/domain"
)

type dtoItem struct {
	domain.Duration
	time.Time
	desc string
}

type dtoList struct {
	domain.Duration
	items     []dtoItem
	startTime time.Time
	finalTime time.Time
}

type service struct {
	domain.Timetable
}

func (s *service) state() dtoList {
	ans := dtoList{items: make([]dtoItem, 0, len(s.Items))}
	for _, i := range s.Timetable.Items {
		ans.Duration = ans.Duration.AddDuration(i.Duration)
		ans.items = append(ans.items, dtoItem{i.Duration, s.Time.Add(ans.TimeDuration()), i.Name})
	}
	ans.startTime = s.Time
	ans.finalTime = s.Time.Add(-ans.TimeDuration())
	return ans
}

func (s *service) settime(t time.Time) {
	s.Time = t
}

func (s *service) size() int {
	return len(s.Items)
}
