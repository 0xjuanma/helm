package tui

import tea "github.com/charmbracelet/bubbletea"

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
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case " ":
		m.pomodoro.Timer.Toggle()
	case "r":
		m.pomodoro.Timer.Reset()
	case "n":
		m.pomodoro.NextPhase()
	}
	return m, nil
}

func (m Model) handleTick() (tea.Model, tea.Cmd) {
	if m.pomodoro.Timer.Tick(tickInterval) {
		m.pomodoro.NextPhase()
	}
	return m, tickCmd()
}
