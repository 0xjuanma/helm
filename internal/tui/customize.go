package tui

import (
	"fmt"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/0xjuanma/helm/internal/config"
)

type editField int

const (
	fieldNone editField = iota
	fieldWorkflowName
	fieldStepName
	fieldStepDuration
)

type editState struct {
	workflowIdx int
	stepIdx     int
	field       editField
	input       string
	draft       *config.WorkflowConfig
	sound       config.SoundConfig
}

func (m Model) updateCustomize(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Sequence(setTitle(""), tea.Quit)
	case "esc":
		m.screen = screenSelect
		return m, nil
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < 1 {
			m.cursor++
		}
	case "enter", " ":
		// 0 = Design (idx 1), 1 = Custom (idx 2)
		workflowIdx := m.cursor + 1
		m.edit = &editState{
			workflowIdx: workflowIdx,
			stepIdx:     0,
			field:       fieldNone,
		}
		m.initDraft()
		m.screen = screenEdit
		m.cursor = 0
	}
	return m, nil
}

func (m *Model) initDraft() {
	var wc *config.WorkflowConfig
	if m.edit.workflowIdx == 1 {
		wc = m.cfg.Design
	} else {
		wc = m.cfg.Custom
	}

	if wc == nil {
		wc = &config.WorkflowConfig{
			Name:           "Custom",
			Steps:          []config.StepConfig{{Name: "STEP 1", Minutes: 10}},
			Loop:           false,
			AutoTransition: true,
		}
	}

	draft := &config.WorkflowConfig{
		Name:           wc.Name,
		Steps:          make([]config.StepConfig, len(wc.Steps)),
		Loop:           wc.Loop,
		AutoTransition: wc.AutoTransition,
	}
	copy(draft.Steps, wc.Steps)
	m.edit.draft = draft

	// Load sound from workflow config, or use default if not set
	if wc.Sound != nil {
		soundCopy := *wc.Sound
		soundCopy.Normalize()
		m.edit.sound = soundCopy
	} else {
		// Use default otherwise
		m.edit.sound = config.DefaultSoundConfig()
	}
}

type editMenuLayout struct {
	soundIdx   int
	stepsStart int
	addIdx     int
	saveIdx    int
	cancelIdx  int
	menuSize   int
}

func (m Model) editMenuLayout() editMenuLayout {
	stepCount := len(m.edit.draft.Steps)
	layout := editMenuLayout{
		soundIdx:   3,
		stepsStart: 4,
	}
	layout.addIdx = layout.stepsStart + stepCount
	layout.saveIdx = layout.addIdx + 1
	layout.cancelIdx = layout.saveIdx + 1
	layout.menuSize = layout.cancelIdx + 1
	return layout
}

func (m Model) updateEdit(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.edit.field != fieldNone {
		return m.updateEditInput(msg)
	}

	layout := m.editMenuLayout()

	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Sequence(setTitle(""), tea.Quit)
	case "esc":
		m.screen = screenCustomize
		m.edit = nil
		m.cursor = 0
		return m, nil
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < layout.menuSize-1 {
			m.cursor++
		}
	case "enter", " ":
		return m.handleEditSelect()
	case "d", "backspace":
		return m.handleDeleteStep()
	}
	return m, nil
}

func (m Model) handleEditSelect() (tea.Model, tea.Cmd) {
	stepCount := len(m.edit.draft.Steps)
	layout := m.editMenuLayout()

	switch {
	case m.cursor == 0:
		m.edit.field = fieldWorkflowName
		m.edit.input = m.edit.draft.Name
	case m.cursor == 1:
		m.edit.draft.Loop = !m.edit.draft.Loop
	case m.cursor == 2:
		m.edit.draft.AutoTransition = !m.edit.draft.AutoTransition
	case m.cursor == layout.soundIdx:
		m.cycleSoundState()
	case m.cursor >= layout.stepsStart && m.cursor < layout.stepsStart+stepCount:
		m.edit.stepIdx = m.cursor - layout.stepsStart
		m.edit.field = fieldStepName
		m.edit.input = m.edit.draft.Steps[m.edit.stepIdx].Name
	case m.cursor == layout.addIdx:
		if stepCount < config.MaxSteps {
			m.edit.draft.Steps = append(m.edit.draft.Steps, config.StepConfig{
				Name:    fmt.Sprintf("STEP %d", stepCount+1),
				Minutes: 10,
			})
		}
	case m.cursor == layout.saveIdx:
		return m.saveWorkflow()
	case m.cursor == layout.cancelIdx:
		m.screen = screenCustomize
		m.edit = nil
		m.cursor = 0
	}
	return m, nil
}

func (m Model) handleDeleteStep() (tea.Model, tea.Cmd) {
	stepCount := len(m.edit.draft.Steps)
	layout := m.editMenuLayout()
	if m.cursor >= layout.stepsStart && m.cursor < layout.stepsStart+stepCount && stepCount > 1 {
		idx := m.cursor - layout.stepsStart
		m.edit.draft.Steps = append(m.edit.draft.Steps[:idx], m.edit.draft.Steps[idx+1:]...)
		if m.cursor >= layout.stepsStart+len(m.edit.draft.Steps) {
			m.cursor--
		}
	}
	return m, nil
}

func (m Model) updateEditInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.edit.field = fieldNone
		m.edit.input = ""
	case "enter":
		return m.applyInput()
	case "backspace":
		if len(m.edit.input) > 0 {
			m.edit.input = m.edit.input[:len(m.edit.input)-1]
		}
	default:
		char := msg.String()
		if len(char) == 1 {
			if m.edit.field == fieldStepDuration {
				if char >= "0" && char <= "9" && len(m.edit.input) < 2 {
					m.edit.input += char
				}
			} else {
				if len(m.edit.input) < 20 {
					m.edit.input += char
				}
			}
		}
	}
	return m, nil
}

func (m Model) applyInput() (tea.Model, tea.Cmd) {
	input := strings.TrimSpace(m.edit.input)
	if input == "" {
		m.edit.field = fieldNone
		m.edit.input = ""
		return m, nil
	}

	switch m.edit.field {
	case fieldWorkflowName:
		m.edit.draft.Name = strings.ToUpper(input)
	case fieldStepName:
		m.edit.draft.Steps[m.edit.stepIdx].Name = strings.ToUpper(input)
		// Move to duration input
		m.edit.field = fieldStepDuration
		m.edit.input = strconv.Itoa(m.edit.draft.Steps[m.edit.stepIdx].Minutes)
		return m, nil
	case fieldStepDuration:
		mins, err := strconv.Atoi(input)
		if err == nil && mins >= config.MinStepMinutes && mins <= config.MaxStepMinutes {
			m.edit.draft.Steps[m.edit.stepIdx].Minutes = mins
		}
	}

	m.edit.field = fieldNone
	m.edit.input = ""
	return m, nil
}

func (m Model) saveWorkflow() (tea.Model, tea.Cmd) {
	if !m.edit.draft.IsValid() {
		return m, nil
	}

	m.edit.sound.Normalize()

	// Always save sound config to workflow (even if unchanged, to ensure it's persisted)
	soundCopy := m.edit.sound
	m.edit.draft.Sound = &soundCopy

	if m.edit.workflowIdx == 1 {
		m.cfg.Design = m.edit.draft
	} else {
		m.cfg.Custom = m.edit.draft
	}

	_ = config.Save(m.cfg)
	m.workflows = m.cfg.BuildWorkflows()

	m.screen = screenCustomize
	m.edit = nil
	m.cursor = 0
	return m, nil
}

func (m Model) viewCustomize() string {
	title := titleStyle.Render("CUSTOMIZE")
	subtitle := subtitleStyle.Render("Select a workflow to edit")

	var items string

	// Get actual workflow names from config
	designName := "Design Interview"
	if m.cfg.Design != nil && m.cfg.Design.Name != "" {
		designName = m.cfg.Design.Name
	}

	customName := "Custom"
	if m.cfg.Custom != nil && m.cfg.Custom.Name != "" {
		customName = m.cfg.Custom.Name
	}

	options := []string{designName, customName}
	for i, name := range options {
		prefix := "  "
		style := itemStyle
		if i == m.cursor {
			prefix = "> "
			style = selectedItemStyle
		}

		var status string
		switch i {
		case 1:
			if m.cfg.Custom == nil {
				status = " [empty]"
			}
		}
		line := fmt.Sprintf("%s%s%s", prefix, name, status)
		items += style.Render(line) + "\n"
	}

	help := helpStyle.Render("[j/k] navigate  [enter] edit  [esc] back  [q] quit")

	return containerStyle.Render(
		lipgloss.JoinVertical(lipgloss.Center, title, subtitle, "", items, help),
	)
}

func (m Model) viewEdit() string {
	title := titleStyle.Render("EDIT WORKFLOW")

	var lines []string
	layout := m.editMenuLayout()

	nameLabel := "Name: "
	nameValue := m.edit.draft.Name
	if m.edit.field == fieldWorkflowName {
		nameValue = m.edit.input + "_"
	}
	lines = append(lines, m.formatEditLine(0, nameLabel+nameValue))

	loopValue := "No"
	if m.edit.draft.Loop {
		loopValue = "Yes"
	}
	lines = append(lines, m.formatEditLine(1, fmt.Sprintf("Loop: %s", loopValue)))

	autoTransValue := "No"
	if m.edit.draft.AutoTransition {
		autoTransValue = "Yes"
	}
	lines = append(lines, m.formatEditLine(2, fmt.Sprintf("Auto-Transition: %s", autoTransValue)))

	soundLabel := "Sound: Off"
	if m.edit.sound.Enabled {
		if m.edit.sound.Mode == config.SoundModeMac {
			soundLabel = "Sound: macOS system sound"
		} else {
			soundLabel = "Sound: Terminal bell"
		}
	}
	lines = append(lines, m.formatEditLine(layout.soundIdx, soundLabel))

	for i, step := range m.edit.draft.Steps {
		stepLine := fmt.Sprintf("  %s (%dm)", step.Name, step.Minutes)
		if m.edit.field == fieldStepName && m.edit.stepIdx == i {
			stepLine = fmt.Sprintf("  Name: %s_", m.edit.input)
		} else if m.edit.field == fieldStepDuration && m.edit.stepIdx == i {
			stepLine = fmt.Sprintf("  Duration: %s_ min", m.edit.input)
		}
		lines = append(lines, m.formatEditLine(layout.stepsStart+i, stepLine))
	}

	addLabel := "+ Add Step"
	if len(m.edit.draft.Steps) >= config.MaxSteps {
		addLabel = "(max steps reached)"
	}
	lines = append(lines, m.formatEditLine(layout.addIdx, addLabel))

	lines = append(lines, "")
	lines = append(lines, m.formatEditLine(layout.saveIdx, "Save"))
	lines = append(lines, m.formatEditLine(layout.cancelIdx, "Cancel"))

	content := strings.Join(lines, "\n")

	var help string
	if m.edit.field != fieldNone {
		help = helpStyle.Render("[enter] confirm  [esc] cancel")
	} else {
		help = helpStyle.Render("[j/k] navigate  [enter] select  [d] delete step  [esc] back")
	}

	return containerStyle.Render(
		lipgloss.JoinVertical(lipgloss.Center, title, "", content, "", help),
	)
}

func (m Model) formatEditLine(idx int, text string) string {
	prefix := "  "
	style := itemStyle
	if m.cursor == idx && m.edit.field == fieldNone {
		prefix = "> "
		style = selectedItemStyle
	}
	return style.Render(prefix + text)
}

func (m *Model) cycleSoundState() {
	// Cycle through: Off → Terminal → Mac → Off
	if !m.edit.sound.Enabled {
		// Off → Terminal
		m.edit.sound.Enabled = true
		m.edit.sound.Mode = config.SoundModeTerminal
	} else if m.edit.sound.Mode == config.SoundModeTerminal {
		// Terminal → Mac
		m.edit.sound.Mode = config.SoundModeMac
		if m.edit.sound.Tone == "" {
			m.edit.sound.Tone = config.DefaultMacTone
		}
	} else {
		// Mac → Off
		m.edit.sound.Enabled = false
	}
}
