package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"

	"github.com/0xjuanma/helm/internal/timer"
)

func (m Model) View() string {
	title := titleStyle.Render(m.pomodoro.Phase.String())
	timerDisplay := m.renderTimer()
	sessions := sessionStyle.Render(fmt.Sprintf("Sessions: %d", m.pomodoro.CompletedWork))
	help := helpStyle.Render("[space] start/pause  [r] reset  [n] skip  [q] quit")

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		timerDisplay,
		sessions,
		help,
	)

	if m.width > 0 && m.height > 0 {
		content = lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
	}

	return containerStyle.Render(content)
}

func (m Model) renderTimer() string {
	remaining := m.pomodoro.Timer.Remaining
	minutes := int(remaining.Minutes())
	seconds := int(remaining.Seconds()) % 60
	display := fmt.Sprintf("%02d:%02d", minutes, seconds)

	if m.pomodoro.Timer.State == timer.Running {
		return timerStyle.Render(display)
	}
	return pausedStyle.Render(display)
}
