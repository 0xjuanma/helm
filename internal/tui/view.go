package tui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"

	"github.com/0xjuanma/helm/internal/timer"
	"github.com/0xjuanma/helm/internal/workflow"
)

func (m Model) View() string {
	var content string
	switch m.screen {
	case screenSelect:
		content = m.viewSelect()
	case screenTimer:
		content = m.viewTimer()
	case screenComplete:
		content = m.viewComplete()
	}

	if m.width > 0 && m.height > 0 {
		content = lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
	}
	return content
}

func (m Model) viewSelect() string {
	title := titleStyle.Render("HELM")
	subtitle := subtitleStyle.Render("Select a workflow")

	var items string
	for i, w := range m.workflows {
		prefix := "  "
		style := itemStyle
		if i == m.cursor {
			prefix = "> "
			style = selectedItemStyle
		}
		totalTime := formatDuration(calcTotalTime(w.Steps))
		loopIndicator := ""
		if w.Loop {
			loopIndicator = " [loop]"
		}
		line := fmt.Sprintf("%s%s (%s)%s", prefix, w.Name, totalTime, loopIndicator)
		items += style.Render(line) + "\n"
	}

	help := helpStyle.Render("[j/k] navigate  [enter] select  [q] quit")

	return containerStyle.Render(
		lipgloss.JoinVertical(lipgloss.Center, title, subtitle, "", items, help),
	)
}

func (m Model) viewTimer() string {
	isRunning := m.session.Timer.State == timer.Running

	// Step label with colored background
	var stepLabel string
	if isRunning {
		stepLabel = stepLabelStyle.Render(m.session.CurrentStepName())
	} else {
		stepLabel = stepLabelPausedStyle.Render(m.session.CurrentStepName())
	}

	// Large timer display
	remaining := m.session.Timer.Remaining
	minutes := int(remaining.Minutes())
	seconds := int(remaining.Seconds()) % 60
	timeStr := fmt.Sprintf("%02d:%02d", minutes, seconds)
	largeTime := renderLargeTime(timeStr)

	var timerDisplay string
	if isRunning {
		timerDisplay = timerLargeStyle.Render(largeTime)
	} else {
		timerDisplay = timerPausedLargeStyle.Render(largeTime)
	}

	// Progress bar
	prog := m.progress()
	var progressDisplay string
	if isRunning {
		progressDisplay = progressContainerStyle.Render(m.progressBar.ViewAs(prog))
	} else {
		pausedBar := newPausedProgressBar()
		progressDisplay = progressContainerStyle.Render(pausedBar.ViewAs(prog))
	}

	// Step progress
	current, total := m.session.StepProgress()
	stepProgress := sessionStyle.Render(fmt.Sprintf("Step %d/%d", current, total))

	// Workflow name
	workflowName := subtitleStyle.Render(m.session.Workflow.Name)

	// Status indicator
	var status string
	if !isRunning {
		status = subtitleStyle.Render("PAUSED")
	}

	help := helpStyle.Render("[space] start/pause  [r] reset  [n] skip  [esc] back  [q] quit")

	elements := []string{workflowName, stepLabel, "", timerDisplay, progressDisplay, stepProgress}
	if status != "" {
		elements = append(elements, status)
	}
	elements = append(elements, help)

	return containerStyle.Render(
		lipgloss.JoinVertical(lipgloss.Center, elements...),
	)
}

func (m Model) viewComplete() string {
	checkmark := renderLargeTime("00:00")
	title := completeStyle.Render("COMPLETE")
	message := subtitleStyle.Render(fmt.Sprintf("%s finished", m.session.Workflow.Name))

	// Full progress bar
	fullBar := progressContainerStyle.Render(m.progressBar.ViewAs(1.0))

	help := helpStyle.Render("[r] restart  [enter] back to menu  [q] quit")

	return containerStyle.Render(
		lipgloss.JoinVertical(lipgloss.Center, title, "", completeStyle.Render(checkmark), fullBar, message, "", help),
	)
}

func calcTotalTime(steps []workflow.Step) time.Duration {
	var total time.Duration
	for _, s := range steps {
		total += s.Duration
	}
	return total
}

func formatDuration(d time.Duration) string {
	m := int(d.Minutes())
	return fmt.Sprintf("%dm", m)
}
