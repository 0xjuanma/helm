package tui

import (
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/lipgloss"
)

const progressWidth = 40

var (
	subtle     = lipgloss.AdaptiveColor{Light: "#666666", Dark: "#888888"}
	highlight  = lipgloss.AdaptiveColor{Light: "#C41E3A", Dark: "#FF6B6B"}
	accent     = lipgloss.AdaptiveColor{Light: "#2E7D32", Dark: "#81C784"}
	muted      = lipgloss.AdaptiveColor{Light: "#AAAAAA", Dark: "#555555"}
	transition = lipgloss.AdaptiveColor{Light: "#FF8C00", Dark: "#FFD700"} // Gold/Orange for transition

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(highlight)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(subtle)

	// Large timer display using ASCII art-style numbers
	timerLargeStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(highlight).
			MarginTop(1).
			MarginBottom(1)

	timerPausedLargeStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(subtle).
				MarginTop(1).
				MarginBottom(1)

	stepLabelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(highlight).
			Padding(0, 2).
			Bold(true)

	stepLabelPausedStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFFFFF")).
				Background(subtle).
				Padding(0, 2).
				Bold(true)

	helpStyle = lipgloss.NewStyle().
			Foreground(muted).
			MarginTop(2)

	sessionStyle = lipgloss.NewStyle().
			Foreground(subtle).
			MarginTop(1)

	itemStyle = lipgloss.NewStyle().
			Foreground(subtle).
			PaddingLeft(2)

	selectedItemStyle = lipgloss.NewStyle().
				Foreground(highlight).
				Bold(true).
				PaddingLeft(0)

	completeStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(accent)

	containerStyle = lipgloss.NewStyle().
			Align(lipgloss.Center).
			Padding(1, 2)

	progressContainerStyle = lipgloss.NewStyle().
				MarginTop(1).
				MarginBottom(1)

	// Transition styles for auto-transition between stages
	transitionStyle = lipgloss.NewStyle().
			Foreground(transition).
			Bold(true)

	transitionPulseStyle = lipgloss.NewStyle().
				Foreground(transition).
				Bold(true)

	transitionDimStyle = lipgloss.NewStyle().
				Foreground(muted).
				Bold(false)

	transitionLabelStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFFFFF")).
				Background(transition).
				Padding(0, 2).
				Bold(true)
)

func newProgressBar() progress.Model {
	p := progress.New(
		progress.WithGradient("#4ECDC4", "#45B7AA"), // Teal gradient
		progress.WithWidth(progressWidth),
		progress.WithoutPercentage(),
	)
	p.Full = '█'
	p.Empty = '░'
	p.EmptyColor = string(muted.Dark)
	return p
}

func newPausedProgressBar() progress.Model {
	p := progress.New(
		progress.WithSolidFill(string(subtle.Dark)),
		progress.WithWidth(progressWidth),
		progress.WithoutPercentage(),
	)
	p.Full = '█'
	p.Empty = '░'
	p.EmptyColor = string(muted.Dark)
	return p
}
