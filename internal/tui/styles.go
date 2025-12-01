package tui

import "github.com/charmbracelet/lipgloss"

var (
	subtle    = lipgloss.AdaptiveColor{Light: "#666666", Dark: "#999999"}
	highlight = lipgloss.AdaptiveColor{Light: "#C41E3A", Dark: "#FF6B6B"}
	accent    = lipgloss.AdaptiveColor{Light: "#2E7D32", Dark: "#81C784"}

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(highlight).
			MarginBottom(1)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(subtle).
			MarginBottom(1)

	timerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(highlight).
			Padding(1, 4)

	pausedStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(subtle).
			Padding(1, 4)

	helpStyle = lipgloss.NewStyle().
			Foreground(subtle).
			MarginTop(1)

	sessionStyle = lipgloss.NewStyle().
			Foreground(accent).
			MarginTop(1)

	itemStyle = lipgloss.NewStyle().
			Foreground(subtle)

	selectedItemStyle = lipgloss.NewStyle().
				Foreground(highlight).
				Bold(true)

	completeStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(accent).
			MarginBottom(1)

	containerStyle = lipgloss.NewStyle().
			Align(lipgloss.Center).
			Padding(2, 4)
)
