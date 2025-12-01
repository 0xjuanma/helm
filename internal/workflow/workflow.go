package workflow

import "time"

type Step struct {
	Name     string
	Duration time.Duration
}

type Workflow struct {
	Name  string
	Steps []Step
	Loop  bool
}

func (w *Workflow) StepCount() int {
	return len(w.Steps)
}
