package workflow

import (
	"testing"
	"time"
)

func TestQuick(t *testing.T) {
	d := 10 * time.Minute
	w := Quick(d)

	if w.Name != "QUICK TIMER" {
		t.Errorf("Name = %q, want %q", w.Name, "QUICK TIMER")
	}
	if w.Loop {
		t.Error("Loop should be false")
	}
	if w.AutoTransition {
		t.Error("AutoTransition should be false")
	}
	if len(w.Steps) != 1 {
		t.Fatalf("len(Steps) = %d, want 1", len(w.Steps))
	}
	if w.Steps[0].Name != "QUICK TIMER" {
		t.Errorf("Steps[0].Name = %q, want %q", w.Steps[0].Name, "QUICK TIMER")
	}
	if w.Steps[0].Duration != d {
		t.Errorf("Steps[0].Duration = %v, want %v", w.Steps[0].Duration, d)
	}
}
