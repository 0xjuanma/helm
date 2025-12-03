package workflow

import "time"

func Pomodoro() Workflow {
	return Workflow{
		Name:           "Pomodoro",
		Loop:           true,
		AutoTransition: true,
		Steps: []Step{
			{Name: "WORK", Duration: 25 * time.Minute},
			{Name: "SHORT BREAK", Duration: 5 * time.Minute},
			{Name: "WORK", Duration: 25 * time.Minute},
			{Name: "SHORT BREAK", Duration: 5 * time.Minute},
			{Name: "WORK", Duration: 25 * time.Minute},
			{Name: "SHORT BREAK", Duration: 5 * time.Minute},
			{Name: "WORK", Duration: 25 * time.Minute},
			{Name: "LONG BREAK", Duration: 15 * time.Minute},
		},
	}
}
