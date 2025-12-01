package timer

import (
	"github.com/0xjuanma/helm/internal/workflow"
)

type Session struct {
	Workflow    *workflow.Workflow
	Timer       *Timer
	CurrentStep int
	Completed   bool
}

func NewSession(w *workflow.Workflow) *Session {
	s := &Session{
		Workflow:    w,
		CurrentStep: 0,
	}
	s.Timer = New(w.Steps[0].Duration)
	return s
}

func (s *Session) CurrentStepName() string {
	return s.Workflow.Steps[s.CurrentStep].Name
}

func (s *Session) StepProgress() (current, total int) {
	return s.CurrentStep + 1, s.Workflow.StepCount()
}

func (s *Session) NextStep() {
	s.CurrentStep++
	if s.CurrentStep >= s.Workflow.StepCount() {
		if s.Workflow.Loop {
			s.CurrentStep = 0
		} else {
			s.Completed = true
			return
		}
	}
	s.Timer = New(s.Workflow.Steps[s.CurrentStep].Duration)
}

func (s *Session) Reset() {
	s.CurrentStep = 0
	s.Completed = false
	s.Timer = New(s.Workflow.Steps[0].Duration)
}

