package tui

import (
	"testing"
	"time"

	"github.com/0xjuanma/helm/internal/config"
	"github.com/0xjuanma/helm/internal/timer"
)

func TestNewQuickModel(t *testing.T) {
	m := NewQuickModel(10)

	if m.screen != screenTimer {
		t.Errorf("screen = %v, want screenTimer", m.screen)
	}
	if m.session == nil {
		t.Fatal("session should not be nil")
	}
	if m.session.Workflow.Name != "QUICK TIMER" {
		t.Errorf("session.Workflow.Name = %q, want %q", m.session.Workflow.Name, "QUICK TIMER")
	}
	if m.session.Timer.Duration != 10*time.Minute {
		t.Errorf("session.Timer.Duration = %v, want %v", m.session.Timer.Duration, 10*time.Minute)
	}
	if m.session.Timer.State != timer.Stopped {
		t.Errorf("session.Timer.State = %v, want Stopped", m.session.Timer.State)
	}
	if m.currentSound != config.DefaultSoundConfig() {
		t.Errorf("currentSound = %v, want %v", m.currentSound, config.DefaultSoundConfig())
	}
	if m.cfg == nil {
		t.Error("cfg should not be nil")
	}
	if len(m.workflows) != 3 {
		t.Errorf("len(workflows) = %d, want 3", len(m.workflows))
	}
}
