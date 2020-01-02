package service

import (
	"time"

	"github.com/waffleboot/timer/domain"
)

//DtoItem ...
type DtoItem struct {
	domain.Duration
	time.Time
	Desc string
}

//DtoList ...
type DtoList struct {
	domain.Duration
	Items     []DtoItem
	StartTime time.Time
	FinalTime time.Time
}

//Service ...
type Service struct {
	domain.Timetable
}

//State ...
func (s *Service) State() DtoList {
	ans := DtoList{Items: make([]DtoItem, 0, len(s.Items))}
	for _, i := range s.Timetable.Items {
		ans.Duration = ans.Duration.AddDuration(i.Duration)
		ans.Items = append(ans.Items, DtoItem{i.Duration, s.Time.Add(ans.TimeDuration()), i.Name})
	}
	ans.StartTime = s.Time
	ans.FinalTime = s.Time.Add(-ans.TimeDuration())
	return ans
}

//SetTime ...
func (s *Service) SetTime(t time.Time) {
	s.Time = t
}

//Size ...
func (s *Service) Size() int {
	return len(s.Items)
}
