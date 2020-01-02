package service_test

import (
	"os"
	"testing"
	"time"

	"github.com/waffleboot/timer/domain"
	"github.com/waffleboot/timer/service"
)

var gTime time.Time

func TestMain(m *testing.M) {
	t, err := time.Parse("15:04", "12:37")
	if err != nil {
		os.Exit(-1)
	}
	gTime = t
	m.Run()
}

func BenchmarkState(b *testing.B) {
	var s service.Service
	for i := 0; i < b.N; i++ {
		s.Add(domain.Duration{Hours: 1, Minutes: 20}, "test")
	}
	s.SetTime(gTime)
	_ = s.State()
}

func TestState(t *testing.T) {
	var s service.Service
	s.Add(domain.Duration{Hours: 1, Minutes: 20}, "")
	s.Add(domain.Duration{Hours: 2, Minutes: 60}, "")
	s.SetTime(gTime)
	state := s.State()
	if state.Duration.Hours != 4 {
		t.Fatal(state.Duration)
	}
	if state.Duration.Minutes != 20 {
		t.Fatal(state.Duration)
	}
	if s.Size() != 2 {
		t.Fatal(s.Size())
	}
}
