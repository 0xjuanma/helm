package tui

import (
	"testing"

	"github.com/0xjuanma/helm/internal/config"
)

func newCustomizeModel(workflowIdx int, cfg *config.Config) Model {
	return Model{
		screen: screenEdit,
		cfg:    cfg,
		edit:   &editState{workflowIdx: workflowIdx},
	}
}

func TestInitDraftRoutesToCorrectSlot(t *testing.T) {
	cfg := &config.Config{
		Custom:  &config.WorkflowConfig{Name: "SLOT 2"},
		Custom2: &config.WorkflowConfig{Name: "SLOT 3"},
		Custom3: &config.WorkflowConfig{Name: "SLOT 4"},
	}

	m := newCustomizeModel(3, cfg)
	m.initDraft()

	if m.edit.draft.Name != "SLOT 3" {
		t.Errorf("initDraft() with workflowIdx=3: draft.Name = %q, want %q (from cfg.Custom2)", m.edit.draft.Name, "SLOT 3")
	}
}

func TestSaveWorkflowRoutesToCorrectSlot(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	cfg := &config.Config{}
	m := newCustomizeModel(4, cfg)
	m.edit.draft = &config.WorkflowConfig{
		Name:  "NEW SLOT 4",
		Steps: []config.StepConfig{{Name: "STEP 1", Minutes: 10}},
	}

	got, _ := m.saveWorkflow()
	gotModel := got.(Model)

	if gotModel.cfg.Custom3 == nil || gotModel.cfg.Custom3.Name != "NEW SLOT 4" {
		t.Errorf("saveWorkflow() with workflowIdx=4: cfg.Custom3 = %+v, want Name %q", gotModel.cfg.Custom3, "NEW SLOT 4")
	}
	if gotModel.cfg.Custom != nil {
		t.Errorf("saveWorkflow() with workflowIdx=4: cfg.Custom = %+v, want nil (should not misroute to slot 2)", gotModel.cfg.Custom)
	}
}
