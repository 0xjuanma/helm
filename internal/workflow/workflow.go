package workflow

import "time"

type Step struct {
	Name     string
	Duration time.Duration
}

type Workflow struct {
	Name  string
	Steps []Step
	Loop  bool // If true, restart from beginning after last step
}

func (w *Workflow) StepCount() int {
	return len(w.Steps)
}

