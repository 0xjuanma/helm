package tui

import (
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/0xjuanma/helm/internal/timer"
	"github.com/0xjuanma/helm/internal/workflow"
)

const tickInterval = time.Second

type screen int

const (
	screenSelect screen = iota
	screenTimer
	screenComplete
)

type tickMsg time.Time

type Model struct {
	screen      screen
	workflows   []workflow.Workflow
	cursor      int
	session     *timer.Session
	progressBar progress.Model
	width       int
	height      int
}

func NewModel() Model {
	return Model{
		screen:      screenSelect,
		workflows:   workflow.Presets,
		cursor:      0,
		progressBar: newProgressBar(),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func tickCmd() tea.Cmd {
	return tea.Tick(tickInterval, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m *Model) progress() float64 {
	if m.session == nil {
		return 0
	}
	elapsed := m.session.Timer.Duration - m.session.Timer.Remaining
	return float64(elapsed) / float64(m.session.Timer.Duration)
}
