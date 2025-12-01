package workflow

import "time"

const MaxWorkflows = 3

var Presets = []Workflow{
	{
		Name: "Pomodoro",
		Loop: true,
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
	},
	{
		Name: "Design Interview",
		Loop: false,
		Steps: []Step{
			{Name: "REQUIREMENTS", Duration: 5 * time.Minute},
			{Name: "ENTITIES & API", Duration: 7 * time.Minute},
			{Name: "HIGH-LEVEL", Duration: 15 * time.Minute},
			{Name: "DEEP-DIVE", Duration: 10 * time.Minute},
		},
	},
	{
		Name: "Quick Focus",
		Loop: true,
		Steps: []Step{
			{Name: "FOCUS", Duration: 15 * time.Minute},
			{Name: "BREAK", Duration: 3 * time.Minute},
		},
	},
}

