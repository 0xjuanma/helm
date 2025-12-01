package tui

import (
	"fmt"

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
	case screenCustomize:
		return m.updateCustomize(msg)
	case screenEdit:
		return m.updateEdit(msg)
	}
	return m, nil
}

func (m Model) handleSelectKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Sequence(setTitle(""), tea.Quit)
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < len(m.workflows)-1 {
			m.cursor++
		}
	case "enter", " ":
		// Don't start empty workflow
		if m.cursor == 2 && m.cfg.Custom == nil {
			return m, nil
		}
		m.session = m.startWorkflow(m.cursor)
		m.screen = screenTimer
		return m, tea.Batch(tickCmd(), m.updateTitle())
	case "c":
		m.screen = screenCustomize
		m.cursor = 0
		return m, nil
	}
	return m, nil
}

func (m Model) handleTimerKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Sequence(setTitle(""), tea.Quit)
	case " ":
		m.session.Timer.Toggle()
		return m, m.updateTitle()
	case "r":
		m.session.Timer.Reset()
		return m, m.updateTitle()
	case "n":
		m.session.NextStep()
		if m.session.Completed {
			m.screen = screenComplete
			return m, tea.Batch(setTitle("Helm - Complete"), bell())
		}
		return m, tea.Batch(m.updateTitle(), bell())
	case "esc":
		m.screen = screenSelect
		m.session = nil
		return m, setTitle("Helm")
	}
	return m, nil
}

func (m Model) handleCompleteKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Sequence(setTitle(""), tea.Quit)
	case "r":
		m.session.Reset()
		m.screen = screenTimer
		return m, tea.Batch(tickCmd(), m.updateTitle())
	case "esc", "enter", " ":
		m.screen = screenSelect
		m.session = nil
		return m, setTitle("Helm")
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
			return m, tea.Batch(setTitle("Helm - Complete"), bell())
		}
		return m, tea.Batch(tickCmd(), m.updateTitle(), bell())
	}
	return m, tea.Batch(tickCmd(), m.updateTitle())
}

func (m Model) startWorkflow(idx int) *timer.Session {
	w := &m.workflows[idx]
	return timer.NewSession(w)
}

func (m Model) updateTitle() tea.Cmd {
	if m.session == nil {
		return setTitle("Helm")
	}
	remaining := m.session.Timer.Remaining
	minutes := int(remaining.Minutes())
	seconds := int(remaining.Seconds()) % 60
	status := ""
	if m.session.Timer.State != timer.Running {
		status = " (paused)"
	}
	title := fmt.Sprintf("%02d:%02d - %s%s", minutes, seconds, m.session.CurrentStepName(), status)
	return setTitle(title)
}

func setTitle(title string) tea.Cmd {
	return func() tea.Msg {
		if title == "" {
			fmt.Print("\033]0;\007")
		} else {
			fmt.Printf("\033]0;%s\007", title)
		}
		return nil
	}
}

func bell() tea.Cmd {
	return func() tea.Msg {
		fmt.Print("\a")
		return nil
	}
}
