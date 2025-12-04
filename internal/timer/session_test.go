package timer

import (
	"testing"
	"time"

	"github.com/0xjuanma/helm/internal/workflow"
)

func newTestWorkflow(loop bool) *workflow.Workflow {
	return &workflow.Workflow{
		Name: "Test",
		Steps: []workflow.Step{
			{Name: "Step 1", Duration: 5 * time.Minute},
			{Name: "Step 2", Duration: 10 * time.Minute},
			{Name: "Step 3", Duration: 3 * time.Minute},
		},
		Loop: loop,
	}
}

func TestNewSession(t *testing.T) {
	w := newTestWorkflow(false)
	s := NewSession(w)

	if s.Workflow != w {
		t.Error("Workflow not set correctly")
	}
	if s.CurrentStep != 0 {
		t.Errorf("CurrentStep = %d, want 0", s.CurrentStep)
	}
	if s.Completed {
		t.Error("Completed should be false")
	}
	if s.Timer.Duration != w.Steps[0].Duration {
		t.Errorf("Timer.Duration = %v, want %v", s.Timer.Duration, w.Steps[0].Duration)
	}
}

func TestCurrentStepName(t *testing.T) {
	w := newTestWorkflow(false)
	s := NewSession(w)

	if got := s.CurrentStepName(); got != "Step 1" {
		t.Errorf("CurrentStepName() = %q, want %q", got, "Step 1")
	}
}

func TestStepProgress(t *testing.T) {
	w := newTestWorkflow(false)
	s := NewSession(w)

	curr, total := s.StepProgress()
	if curr != 1 || total != 3 {
		t.Errorf("StepProgress() = (%d, %d), want (1, 3)", curr, total)
	}
}

func TestNextStep(t *testing.T) {
	w := newTestWorkflow(false)
	s := NewSession(w)

	s.NextStep()

	if s.CurrentStep != 1 {
		t.Errorf("CurrentStep = %d, want 1", s.CurrentStep)
	}
	if s.Timer.Duration != w.Steps[1].Duration {
		t.Errorf("Timer.Duration = %v, want %v", s.Timer.Duration, w.Steps[1].Duration)
	}
	if s.Completed {
		t.Error("Completed should be false")
	}
}

func TestNextStepLoop(t *testing.T) {
	w := newTestWorkflow(true)
	s := NewSession(w)
	s.CurrentStep = 2

	s.NextStep()

	if s.CurrentStep != 0 {
		t.Errorf("CurrentStep = %d, want 0 (should loop)", s.CurrentStep)
	}
	if s.Completed {
		t.Error("Completed should be false when looping")
	}
	if s.Timer.Duration != w.Steps[0].Duration {
		t.Errorf("Timer.Duration = %v, want %v", s.Timer.Duration, w.Steps[0].Duration)
	}
}

func TestNextStepComplete(t *testing.T) {
	w := newTestWorkflow(false)
	s := NewSession(w)
	s.CurrentStep = 2

	s.NextStep()

	if !s.Completed {
		t.Error("Completed should be true when no loop and at end")
	}
	if s.CurrentStep != 3 {
		t.Errorf("CurrentStep = %d, want 3", s.CurrentStep)
	}
}

func TestSessionReset(t *testing.T) {
	w := newTestWorkflow(false)
	s := NewSession(w)
	s.CurrentStep = 2
	s.Completed = true
	s.Timer.Remaining = 0

	s.Reset()

	if s.CurrentStep != 0 {
		t.Errorf("CurrentStep = %d, want 0", s.CurrentStep)
	}
	if s.Completed {
		t.Error("Completed should be false after reset")
	}
	if s.Timer.Duration != w.Steps[0].Duration {
		t.Errorf("Timer.Duration = %v, want %v", s.Timer.Duration, w.Steps[0].Duration)
	}
	if s.Timer.Remaining != w.Steps[0].Duration {
		t.Errorf("Timer.Remaining = %v, want %v", s.Timer.Remaining, w.Steps[0].Duration)
	}
}

func TestNextStepOnCompletedSession(t *testing.T) {
	w := newTestWorkflow(false)
	s := NewSession(w)
	s.CurrentStep = 2
	s.NextStep() // completes the session

	if !s.Completed {
		t.Fatal("Session should be completed")
	}

	// calling NextStep on completed session should be a no-op
	prevStep := s.CurrentStep
	s.NextStep()

	if s.CurrentStep != prevStep {
		t.Errorf("CurrentStep changed from %d to %d, should remain unchanged", prevStep, s.CurrentStep)
	}
}

func TestNewSessionWithEmptyWorkflow(t *testing.T) {
	w := &workflow.Workflow{
		Name:  "Empty",
		Steps: []workflow.Step{},
	}

	s := NewSession(w)

	if s != nil {
		t.Error("NewSession should return nil for empty workflow")
	}
}
