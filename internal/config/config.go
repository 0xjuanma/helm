package config

import (
	"time"

	"github.com/0xjuanma/helm/internal/workflow"
)

const (
	MaxSteps       = 10
	MaxStepMinutes = 180
	MinStepMinutes = 1
)

type StepConfig struct {
	Name    string `json:"name"`
	Minutes int    `json:"minutes"`
}

type WorkflowConfig struct {
	Name           string       `json:"name"`
	Steps          []StepConfig `json:"steps"`
	Loop           bool         `json:"loop"`
	AutoTransition bool         `json:"auto_transition"`
	Sound          *SoundConfig `json:"sound,omitempty"`
}

type SoundMode string

const (
	SoundModeTerminal SoundMode = "terminal"
	SoundModeMac      SoundMode = "mac"
)

type SoundConfig struct {
	Enabled bool      `json:"enabled"`
	Mode    SoundMode `json:"mode"`
	Tone    string    `json:"tone"`
}

const DefaultMacTone = "Ping"

type Config struct {
	Design             *WorkflowConfig `json:"design,omitempty"`
	Custom             *WorkflowConfig `json:"custom,omitempty"`
	Custom2            *WorkflowConfig `json:"custom2,omitempty"`
	Custom3            *WorkflowConfig `json:"custom3,omitempty"`
	TransitionDelaySec int             `json:"transition_delay_sec"` // Delay in seconds before next stage starts (1-10)
}

const (
	DefaultTransitionDelay = 3
	MinTransitionDelay     = 1
	MaxTransitionDelay     = 10
)

func DefaultSoundConfig() SoundConfig {
	return SoundConfig{
		Enabled: true,
		Mode:    SoundModeTerminal,
		Tone:    DefaultMacTone,
	}
}

func DefaultConfig() *Config {
	return &Config{
		Design: &WorkflowConfig{
			Name:           "Design Interview",
			Loop:           false,
			AutoTransition: true,
			Steps: []StepConfig{
				{Name: "REQUIREMENTS", Minutes: 5},
				{Name: "ENTITIES & API", Minutes: 7},
				{Name: "HIGH-LEVEL", Minutes: 15},
				{Name: "DEEP-DIVE", Minutes: 10},
			},
		},
		Custom:             nil,
		Custom2:            nil,
		Custom3:            nil,
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

func (cfg *Config) Normalize() {
	// Normalize sound configs in workflows
	if cfg.Design != nil && cfg.Design.Sound != nil {
		cfg.Design.Sound.Normalize()
	}
	if cfg.Custom != nil && cfg.Custom.Sound != nil {
		cfg.Custom.Sound.Normalize()
	}
	if cfg.Custom2 != nil && cfg.Custom2.Sound != nil {
		cfg.Custom2.Sound.Normalize()
	}
	if cfg.Custom3 != nil && cfg.Custom3.Sound != nil {
		cfg.Custom3.Sound.Normalize()
	}
}

// WorkflowConfigAt returns the *WorkflowConfig backing customizable slot idx (1-4), or nil for
// an out-of-range idx (including the immutable Pomodoro slot, 0).
func (cfg *Config) WorkflowConfigAt(idx int) *WorkflowConfig {
	switch idx {
	case 1:
		return cfg.Design
	case 2:
		return cfg.Custom
	case 3:
		return cfg.Custom2
	case 4:
		return cfg.Custom3
	default:
		return nil
	}
}

// SetWorkflowConfigAt sets the *WorkflowConfig backing customizable slot idx (1-4). No-op for
// an out-of-range idx.
func (cfg *Config) SetWorkflowConfigAt(idx int, wc *WorkflowConfig) {
	switch idx {
	case 1:
		cfg.Design = wc
	case 2:
		cfg.Custom = wc
	case 3:
		cfg.Custom2 = wc
	case 4:
		cfg.Custom3 = wc
	}
}

// GetWorkflowSound returns the sound config for a workflow by index, or default if not set
func (cfg *Config) GetWorkflowSound(idx int) SoundConfig {
	wc := cfg.WorkflowConfigAt(idx)
	if wc != nil && wc.Sound != nil {
		sound := *wc.Sound
		sound.Normalize()
		return sound
	}
	return DefaultSoundConfig()
}

func (cfg *Config) BuildWorkflows() []workflow.Workflow {
	workflows := make([]workflow.Workflow, 5)

	// Slot 0: Pomodoro (immutable)
	workflows[0] = workflow.Pomodoro()

	// Slot 1: Design (customizable)
	if cfg.Design != nil {
		workflows[1] = cfg.Design.ToWorkflow()
	} else {
		workflows[1] = DefaultConfig().Design.ToWorkflow()
	}

	// Slots 2-4: Custom, Custom2, Custom3 (user-created)
	for idx := 2; idx <= 4; idx++ {
		if wc := cfg.WorkflowConfigAt(idx); wc != nil {
			workflows[idx] = wc.ToWorkflow()
		} else {
			workflows[idx] = workflow.Workflow{Name: "Empty - press [c] to customize", AutoTransition: true}
		}
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
		Name:           wc.Name,
		Steps:          steps,
		Loop:           wc.Loop,
		AutoTransition: wc.AutoTransition,
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
		Name:           w.Name,
		Steps:          steps,
		Loop:           w.Loop,
		AutoTransition: w.AutoTransition,
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

func (sc *SoundConfig) Normalize() {
	if sc.Mode == "" && !sc.Enabled {
		sc.Enabled = true
		sc.Mode = SoundModeTerminal
		return
	}

	if sc.Mode != SoundModeTerminal && sc.Mode != SoundModeMac {
		sc.Mode = SoundModeTerminal
	}

	if sc.Tone == "" {
		sc.Tone = DefaultMacTone
	}
}
