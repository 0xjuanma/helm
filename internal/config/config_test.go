package config

import "testing"

func TestBuildWorkflowsHasFiveSlots(t *testing.T) {
	cfg := DefaultConfig()
	workflows := cfg.BuildWorkflows()

	if len(workflows) != 5 {
		t.Fatalf("BuildWorkflows() returned %d slots, want 5", len(workflows))
	}
}

func TestWorkflowConfigAtRoundTrip(t *testing.T) {
	cfg := &Config{}
	wc := &WorkflowConfig{Name: "TEST"}

	for idx := 1; idx <= 4; idx++ {
		cfg.SetWorkflowConfigAt(idx, wc)
		if got := cfg.WorkflowConfigAt(idx); got != wc {
			t.Errorf("slot %d: WorkflowConfigAt() = %v, want %v", idx, got, wc)
		}
		cfg.SetWorkflowConfigAt(idx, nil)
		if got := cfg.WorkflowConfigAt(idx); got != nil {
			t.Errorf("slot %d: WorkflowConfigAt() after clear = %v, want nil", idx, got)
		}
	}
}

func TestWorkflowConfigAtOutOfRange(t *testing.T) {
	cfg := DefaultConfig()
	for _, idx := range []int{0, -1, 5} {
		if got := cfg.WorkflowConfigAt(idx); got != nil {
			t.Errorf("WorkflowConfigAt(%d) = %v, want nil", idx, got)
		}
	}
}
