package config

import "testing"

func TestSaveLoadRoundTripsNewSlots(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	cfg := DefaultConfig()
	cfg.Custom2 = &WorkflowConfig{Name: "CUSTOM TIMER 2", Steps: []StepConfig{{Name: "STEP 1", Minutes: 10}}}
	cfg.Custom3 = &WorkflowConfig{Name: "CUSTOM TIMER 3", Steps: []StepConfig{{Name: "STEP 1", Minutes: 20}}}

	if err := Save(cfg); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	loaded, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if loaded.Custom2 == nil || loaded.Custom2.Name != "CUSTOM TIMER 2" {
		t.Errorf("loaded.Custom2 = %+v, want Name %q", loaded.Custom2, "CUSTOM TIMER 2")
	}
	if loaded.Custom3 == nil || loaded.Custom3.Name != "CUSTOM TIMER 3" {
		t.Errorf("loaded.Custom3 = %+v, want Name %q", loaded.Custom3, "CUSTOM TIMER 3")
	}
}
