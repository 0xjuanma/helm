package timer

import "time"

type State int

const (
	Stopped State = iota
	Running
	Paused
)

type Timer struct {
	Duration  time.Duration
	Remaining time.Duration
	State     State
}

func New(d time.Duration) *Timer {
	return &Timer{
		Duration:  d,
		Remaining: d,
		State:     Stopped,
	}
}

func (t *Timer) Start() {
	t.State = Running
}

func (t *Timer) Pause() {
	if t.State == Running {
		t.State = Paused
	}
}

func (t *Timer) Toggle() {
	switch t.State {
	case Running:
		t.Pause()
	case Stopped, Paused:
		t.Start()
	}
}

func (t *Timer) Reset() {
	t.ResetWith(t.Duration)
}

func (t *Timer) ResetWith(d time.Duration) {
	t.Duration = d
	t.Remaining = d
	t.State = Stopped
}

func (t *Timer) Tick(d time.Duration) bool {
	if t.State != Running {
		return false
	}
	t.Remaining -= d
	if t.Remaining <= 0 {
		t.Remaining = 0
		t.State = Stopped
		return true
	}
	return false
}

func (t *Timer) IsRunning() bool {
	return t.State == Running
}

func (t *Timer) IsComplete() bool {
	return t.Remaining <= 0
}
