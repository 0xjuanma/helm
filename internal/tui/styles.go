package tui

import (
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/lipgloss"
)

const progressWidth = 40

var (
	subtle    = lipgloss.AdaptiveColor{Light: "#666666", Dark: "#888888"}
	highlight = lipgloss.AdaptiveColor{Light: "#3D7A9E", Dark: "#7EB5D6"} // Steel blue (from helm logo)
	gold      = lipgloss.AdaptiveColor{Light: "#C7920A", Dark: "#F9A825"} // Gold (from helm logo)
	muted     = lipgloss.AdaptiveColor{Light: "#AAAAAA", Dark: "#555555"}
	white     = lipgloss.AdaptiveColor{Light: "#333333", Dark: "#FFFFFF"}

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(highlight)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(subtle)

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
			Foreground(lipgloss.Color("#1a1a1a")). // Dark text on gold
			Background(gold).
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
			Foreground(gold)

	completeTimerStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(white)

	containerStyle = lipgloss.NewStyle().
			Align(lipgloss.Center).
			Padding(1, 2)

	progressContainerStyle = lipgloss.NewStyle().
				MarginTop(1).
				MarginBottom(1)

	transitionStyle = lipgloss.NewStyle().
			Foreground(gold).
			Bold(true)

	transitionPulseStyle = lipgloss.NewStyle().
				Foreground(gold).
				Bold(true)

	transitionDimStyle = lipgloss.NewStyle().
				Foreground(muted).
				Bold(false)

	transitionLabelStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#1a1a1a")). // Dark text on gold
				Background(gold).
				Padding(0, 2).
				Bold(true)
)

func newProgressBar() progress.Model {
	p := progress.New(
		progress.WithGradient("#7EB5D6", "#F9A825"), // Steel → Gold gradient
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
