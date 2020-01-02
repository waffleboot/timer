package domain_test

import "testing"

import "github.com/waffleboot/timer/domain"

func TestDuration(t *testing.T) {
	a := domain.Duration{1, 20}
	b := domain.Duration{2, 60}
	c := a.AddDuration(b)
	if c.Hours != 4 {
		t.Fatal(c.Hours)
	}
	if c.Minutes != 20 {
		t.FailNow()
	}
}
