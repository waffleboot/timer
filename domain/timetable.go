package domain

import "time"

type item struct {
	Name string
	Duration
}

//Timetable ...
type Timetable struct {
	Items []item
	time.Time
}

//Add ...
func (t *Timetable) Add(d Duration, name string) {
	if d.isValid() {
		t.Items = append(t.Items, item{name, d})
	}
}

//Del ...
func (t *Timetable) Del(p int) {
	t.Items = append(t.Items[:p-1], t.Items[p:]...)
}

//Clear ...
func (t *Timetable) Clear() {
	t.Items = t.Items[:0]
}
