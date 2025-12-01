package timer

import "time"

type Phase int

const (
	Work Phase = iota
	ShortBreak
	LongBreak
)

const (
	WorkDuration       = 25 * time.Minute
	ShortBreakDuration = 5 * time.Minute
	LongBreakDuration  = 15 * time.Minute
	SessionsUntilLong  = 4
)

func (p Phase) String() string {
	switch p {
	case Work:
		return "WORK"
	case ShortBreak:
		return "SHORT BREAK"
	case LongBreak:
		return "LONG BREAK"
	default:
		return "UNKNOWN"
	}
}

func (p Phase) Duration() time.Duration {
	switch p {
	case Work:
		return WorkDuration
	case ShortBreak:
		return ShortBreakDuration
	case LongBreak:
		return LongBreakDuration
	default:
		return WorkDuration
	}
}

type Pomodoro struct {
	Timer         *Timer
	Phase         Phase
	CompletedWork int
}

func NewPomodoro() *Pomodoro {
	return &Pomodoro{
		Timer: New(WorkDuration),
		Phase: Work,
	}
}

func (p *Pomodoro) NextPhase() {
	if p.Phase == Work {
		p.CompletedWork++
		if p.CompletedWork%SessionsUntilLong == 0 {
			p.Phase = LongBreak
		} else {
			p.Phase = ShortBreak
		}
	} else {
		p.Phase = Work
	}
	p.Timer = New(p.Phase.Duration())
}
