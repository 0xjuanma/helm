package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/0xjuanma/helm/internal/timer"
)

const tickInterval = time.Second

type tickMsg time.Time

type Model struct {
	pomodoro *timer.Pomodoro
	width    int
	height   int
}

func NewModel() Model {
	return Model{
		pomodoro: timer.NewPomodoro(),
	}
}

func (m Model) Init() tea.Cmd {
	return tickCmd()
}

func tickCmd() tea.Cmd {
	return tea.Tick(tickInterval, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
