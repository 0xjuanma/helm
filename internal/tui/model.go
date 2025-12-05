package tui

import (
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/0xjuanma/helm/internal/config"
	"github.com/0xjuanma/helm/internal/timer"
	"github.com/0xjuanma/helm/internal/workflow"
)

const tickInterval = time.Second

type screen int

const (
	screenSelect screen = iota
	screenTimer
	screenComplete
	screenCustomize
	screenEdit
	screenSound
)

type tickMsg time.Time

type Model struct {
	screen          screen
	workflows       []workflow.Workflow
	cfg             *config.Config
	cursor          int
	session         *timer.Session
	progressBar     progress.Model
	edit            *editState
	width           int
	height          int
	transitioning   bool
	transitionTicks int
}

func NewModel() Model {
	cfg, _ := config.Load()
	return Model{
		screen:      screenSelect,
		workflows:   cfg.BuildWorkflows(),
		cfg:         cfg,
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
