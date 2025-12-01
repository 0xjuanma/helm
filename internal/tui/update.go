package tui

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/0xjuanma/helm/internal/timer"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKey(msg)
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tickMsg:
		return m.handleTick()
	}
	return m, nil
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.screen {
	case screenSelect:
		return m.handleSelectKey(msg)
	case screenTimer:
		return m.handleTimerKey(msg)
	case screenComplete:
		return m.handleCompleteKey(msg)
	}
	return m, nil
}

func (m Model) handleSelectKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < len(m.workflows)-1 {
			m.cursor++
		}
	case "enter", " ":
		m.session = m.startWorkflow(m.cursor)
		m.screen = screenTimer
		return m, tickCmd()
	}
	return m, nil
}

func (m Model) handleTimerKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case " ":
		m.session.Timer.Toggle()
	case "r":
		m.session.Timer.Reset()
	case "n":
		m.session.NextStep()
		if m.session.Completed {
			m.screen = screenComplete
		}
	case "esc":
		m.screen = screenSelect
		m.session = nil
		return m, nil
	}
	return m, nil
}

func (m Model) handleCompleteKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "r":
		m.session.Reset()
		m.screen = screenTimer
		return m, tickCmd()
	case "esc", "enter", " ":
		m.screen = screenSelect
		m.session = nil
	}
	return m, nil
}

func (m Model) handleTick() (tea.Model, tea.Cmd) {
	if m.screen != screenTimer || m.session == nil {
		return m, nil
	}
	if m.session.Timer.Tick(tickInterval) {
		m.session.NextStep()
		if m.session.Completed {
			m.screen = screenComplete
			return m, nil
		}
	}
	return m, tickCmd()
}

func (m Model) startWorkflow(idx int) *timer.Session {
	w := &m.workflows[idx]
	return timer.NewSession(w)
}
