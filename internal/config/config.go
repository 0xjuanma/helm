package config

import (
	"time"

	"github.com/0xjuanma/helm/internal/workflow"
)

const (
	MaxSteps       = 10
	MaxStepMinutes = 60
	MinStepMinutes = 1
)

type StepConfig struct {
	Name    string `json:"name"`
	Minutes int    `json:"minutes"`
}

type WorkflowConfig struct {
	Name  string       `json:"name"`
	Steps []StepConfig `json:"steps"`
	Loop  bool         `json:"loop"`
}

type Config struct {
	Design             *WorkflowConfig `json:"design,omitempty"`
	Custom             *WorkflowConfig `json:"custom,omitempty"`
	AutoTransition     bool            `json:"auto_transition"`      // Enable auto-transition between stages
	TransitionDelaySec int             `json:"transition_delay_sec"` // Delay in seconds before next stage starts (1-10)
}

const (
	DefaultTransitionDelay = 3
	MinTransitionDelay     = 1
	MaxTransitionDelay     = 10
)

func DefaultConfig() *Config {
	return &Config{
		Design: &WorkflowConfig{
			Name: "Design Interview",
			Loop: false,
			Steps: []StepConfig{
				{Name: "REQUIREMENTS", Minutes: 5},
				{Name: "ENTITIES & API", Minutes: 7},
				{Name: "HIGH-LEVEL", Minutes: 15},
				{Name: "DEEP-DIVE", Minutes: 10},
			},
		},
		Custom:             nil,
		AutoTransition:     true,
		TransitionDelaySec: DefaultTransitionDelay,
	}
}

// GetTransitionDelay returns the transition delay, clamped to valid range
func (cfg *Config) GetTransitionDelay() int {
	if cfg.TransitionDelaySec < MinTransitionDelay {
		return DefaultTransitionDelay
	}
	if cfg.TransitionDelaySec > MaxTransitionDelay {
		return MaxTransitionDelay
	}
	return cfg.TransitionDelaySec
}

func (cfg *Config) BuildWorkflows() []workflow.Workflow {
	workflows := make([]workflow.Workflow, 3)

	// Slot 0: Pomodoro (immutable)
	workflows[0] = workflow.Pomodoro()

	// Slot 1: Design (customizable)
	if cfg.Design != nil {
		workflows[1] = cfg.Design.ToWorkflow()
	} else {
		workflows[1] = DefaultConfig().Design.ToWorkflow()
	}

	// Slot 2: Custom (user-created)
	if cfg.Custom != nil {
		workflows[2] = cfg.Custom.ToWorkflow()
	} else {
		workflows[2] = workflow.Workflow{Name: "Empty - press [c] to customize"}
	}

	return workflows
}

func (wc *WorkflowConfig) ToWorkflow() workflow.Workflow {
	steps := make([]workflow.Step, len(wc.Steps))
	for i, s := range wc.Steps {
		steps[i] = workflow.Step{
			Name:     s.Name,
			Duration: time.Duration(s.Minutes) * time.Minute,
		}
	}
	return workflow.Workflow{
		Name:  wc.Name,
		Steps: steps,
		Loop:  wc.Loop,
	}
}

func FromWorkflow(w *workflow.Workflow) *WorkflowConfig {
	steps := make([]StepConfig, len(w.Steps))
	for i, s := range w.Steps {
		steps[i] = StepConfig{
			Name:    s.Name,
			Minutes: int(s.Duration.Minutes()),
		}
	}
	return &WorkflowConfig{
		Name:  w.Name,
		Steps: steps,
		Loop:  w.Loop,
	}
}

func (wc *WorkflowConfig) IsValid() bool {
	if wc == nil || wc.Name == "" || len(wc.Steps) == 0 {
		return false
	}
	if len(wc.Steps) > MaxSteps {
		return false
	}
	for _, s := range wc.Steps {
		if s.Name == "" || s.Minutes < MinStepMinutes || s.Minutes > MaxStepMinutes {
			return false
		}
	}
	return true
}
