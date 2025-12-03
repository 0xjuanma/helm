package tui

import (
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/lipgloss"
)

const progressWidth = 40

var (
	subtle    = lipgloss.AdaptiveColor{Light: "#666666", Dark: "#888888"}
	highlight = lipgloss.AdaptiveColor{Light: "#C41E3A", Dark: "#FF6B6B"}
	teal      = lipgloss.AdaptiveColor{Light: "#3BA99C", Dark: "#4ECDC4"}
	accent    = lipgloss.AdaptiveColor{Light: "#2E7D32", Dark: "#81C784"}
	muted     = lipgloss.AdaptiveColor{Light: "#AAAAAA", Dark: "#555555"}
	white     = lipgloss.AdaptiveColor{Light: "#333333", Dark: "#FFFFFF"}

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(highlight)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(subtle)

	// Large timer display using ASCII art-style numbers
	timerLargeStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(highlight). // Red accent
			MarginTop(1).
			MarginBottom(1)

	timerPausedLargeStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(subtle).
				MarginTop(1).
				MarginBottom(1)

	stepLabelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#1a1a1a")). // Dark text on teal
			Background(teal).
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
			Foreground(teal)

	completeTimerStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(white)

	containerStyle = lipgloss.NewStyle().
			Align(lipgloss.Center).
			Padding(1, 2)

	progressContainerStyle = lipgloss.NewStyle().
				MarginTop(1).
				MarginBottom(1)

	// Transition styles for auto-transition between stages
	transitionStyle = lipgloss.NewStyle().
			Foreground(teal).
			Bold(true)

	transitionPulseStyle = lipgloss.NewStyle().
				Foreground(teal).
				Bold(true)

	transitionDimStyle = lipgloss.NewStyle().
				Foreground(muted).
				Bold(false)

	transitionLabelStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#1a1a1a")). // Dark text on teal
				Background(teal).
				Padding(0, 2).
				Bold(true)
)

func newProgressBar() progress.Model {
	p := progress.New(
		progress.WithGradient("#FF6B6B", "#4ECDC4"), // Red → Teal gradient
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
