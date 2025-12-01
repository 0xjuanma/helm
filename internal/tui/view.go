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
		cursor := "  "
		style := itemStyle
		if i == m.cursor {
			cursor = "> "
			style = selectedItemStyle
		}
		totalTime := formatDuration(calcTotalTime(w.Steps))
		loopIndicator := ""
		if w.Loop {
			loopIndicator = " [loop]"
		}
		line := fmt.Sprintf("%s%s (%s)%s", cursor, w.Name, totalTime, loopIndicator)
		items += style.Render(line) + "\n"
	}

	help := helpStyle.Render("[j/k] navigate  [enter] select  [q] quit")

	return containerStyle.Render(
		lipgloss.JoinVertical(lipgloss.Center, title, subtitle, "", items, help),
	)
}

func (m Model) viewTimer() string {
	title := titleStyle.Render(m.session.CurrentStepName())

	current, total := m.session.StepProgress()
	progress := sessionStyle.Render(fmt.Sprintf("Step %d of %d", current, total))

	timerDisplay := m.renderTimer()

	workflowName := subtitleStyle.Render(m.session.Workflow.Name)
	help := helpStyle.Render("[space] start/pause  [r] reset  [n] skip  [esc] back  [q] quit")

	return containerStyle.Render(
		lipgloss.JoinVertical(lipgloss.Center, workflowName, title, timerDisplay, progress, help),
	)
}

func (m Model) viewComplete() string {
	title := completeStyle.Render("COMPLETE")
	message := subtitleStyle.Render(fmt.Sprintf("%s finished", m.session.Workflow.Name))
	help := helpStyle.Render("[r] restart  [enter] back to menu  [q] quit")

	return containerStyle.Render(
		lipgloss.JoinVertical(lipgloss.Center, title, "", message, "", help),
	)
}

func (m Model) renderTimer() string {
	remaining := m.session.Timer.Remaining
	minutes := int(remaining.Minutes())
	seconds := int(remaining.Seconds()) % 60
	display := fmt.Sprintf("%02d:%02d", minutes, seconds)

	if m.session.Timer.State == timer.Running {
		return timerStyle.Render(display)
	}
	return pausedStyle.Render(display)
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
