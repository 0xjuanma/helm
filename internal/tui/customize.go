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

// Sound menu indices are computed dynamically based on mode; see soundMenuItems.

type editState struct {
	workflowIdx int
	stepIdx     int
	field       editField
	input       string
	draft       *config.WorkflowConfig
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
		if m.cursor < 2 {
			m.cursor++
		}
	case "enter", " ":
		if m.cursor == 2 {
			m.screen = screenSound
			m.cursor = 0
			return m, nil
		}

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
}

func (m Model) updateEdit(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.edit.field != fieldNone {
		return m.updateEditInput(msg)
	}

	stepCount := len(m.edit.draft.Steps)
	// Menu: Name, Loop, Auto-Transition, [Steps...], Add Step, Save, Cancel
	menuSize := 4 + stepCount + 1

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
		if m.cursor < menuSize-1 {
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

	switch {
	case m.cursor == 0:
		m.edit.field = fieldWorkflowName
		m.edit.input = m.edit.draft.Name
	case m.cursor == 1:
		m.edit.draft.Loop = !m.edit.draft.Loop
	case m.cursor == 2:
		m.edit.draft.AutoTransition = !m.edit.draft.AutoTransition
	case m.cursor >= 3 && m.cursor < 3+stepCount:
		m.edit.stepIdx = m.cursor - 3
		m.edit.field = fieldStepName
		m.edit.input = m.edit.draft.Steps[m.edit.stepIdx].Name
	case m.cursor == 3+stepCount:
		if stepCount < config.MaxSteps {
			m.edit.draft.Steps = append(m.edit.draft.Steps, config.StepConfig{
				Name:    fmt.Sprintf("STEP %d", stepCount+1),
				Minutes: 10,
			})
		}
	case m.cursor == 4+stepCount:
		return m.saveWorkflow()
	case m.cursor == 5+stepCount:
		m.screen = screenCustomize
		m.edit = nil
		m.cursor = 0
	}
	return m, nil
}

func (m Model) handleDeleteStep() (tea.Model, tea.Cmd) {
	stepCount := len(m.edit.draft.Steps)
	if m.cursor >= 3 && m.cursor < 3+stepCount && stepCount > 1 {
		idx := m.cursor - 3
		m.edit.draft.Steps = append(m.edit.draft.Steps[:idx], m.edit.draft.Steps[idx+1:]...)
		if m.cursor >= 3+len(m.edit.draft.Steps) {
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
	subtitle := subtitleStyle.Render("Select workflow or sound to edit")

	var items string
	options := []string{"Design Interview", "Custom", "Sound"}
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
		case 2:
			mode := "bell"
			if m.cfg.Sound.Mode == config.SoundModeMac {
				mode = "mac"
			}
			if !m.cfg.Sound.Enabled {
				status = " [off]"
			} else {
				status = fmt.Sprintf(" [%s]", mode)
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

	for i, step := range m.edit.draft.Steps {
		stepLine := fmt.Sprintf("  %s (%dm)", step.Name, step.Minutes)
		if m.edit.field == fieldStepName && m.edit.stepIdx == i {
			stepLine = fmt.Sprintf("  Name: %s_", m.edit.input)
		} else if m.edit.field == fieldStepDuration && m.edit.stepIdx == i {
			stepLine = fmt.Sprintf("  Duration: %s_ min", m.edit.input)
		}
		lines = append(lines, m.formatEditLine(3+i, stepLine))
	}

	addIdx := 3 + len(m.edit.draft.Steps)
	addLabel := "+ Add Step"
	if len(m.edit.draft.Steps) >= config.MaxSteps {
		addLabel = "(max steps reached)"
	}
	lines = append(lines, m.formatEditLine(addIdx, addLabel))

	lines = append(lines, "")
	lines = append(lines, m.formatEditLine(addIdx+1, "Save"))
	lines = append(lines, m.formatEditLine(addIdx+2, "Cancel"))

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

func (m Model) updateSound(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	items := m.soundMenuItems()
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Sequence(setTitle(""), tea.Quit)
	case "esc":
		m.screen = screenCustomize
		m.cursor = 2
		return m, nil
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < len(items)-1 {
			m.cursor++
		}
	case "enter", " ":
		return m.handleSoundSelect()
	}
	return m, nil
}

func (m Model) handleSoundSelect() (tea.Model, tea.Cmd) {
	mode := m.cfg.Sound.Mode
	macMode := mode == config.SoundModeMac
	toneIdx := -1
	testIdx := 2
	backIdx := 3
	if macMode {
		toneIdx = 2
		testIdx = 3
		backIdx = 4
	}

	switch {
	case m.cursor == 0:
		m.cfg.Sound.Enabled = !m.cfg.Sound.Enabled
	case m.cursor == 1:
		m.cfg.Sound.Mode = nextSoundMode(m.cfg.Sound.Mode)
		if m.cfg.Sound.Mode == config.SoundModeMac && m.cfg.Sound.Tone == "" {
			m.cfg.Sound.Tone = config.DefaultMacTone
		}
	case macMode && m.cursor == toneIdx:
		m.cfg.Sound.Tone = nextMacTone(m.cfg.Sound.Tone)
	case m.cursor == testIdx:
		return m, bell(m.cfg.Sound)
	case m.cursor == backIdx:
		m.screen = screenCustomize
		m.cursor = 2
		return m, nil
	}

	m.cfg.Sound.Normalize()
	_ = config.Save(m.cfg)
	return m, nil
}

func nextSoundMode(current config.SoundMode) config.SoundMode {
	switch current {
	case config.SoundModeTerminal:
		return config.SoundModeMac
	default:
		return config.SoundModeTerminal
	}
}

func nextMacTone(current string) string {
	tone := current
	if tone == "" {
		tone = config.DefaultMacTone
	}
	tones := macTones()
	for i, t := range tones {
		if strings.EqualFold(t, tone) {
			return tones[(i+1)%len(tones)]
		}
	}
	return config.DefaultMacTone
}

func (m Model) viewSound() string {
	title := titleStyle.Render("SOUND")
	subtitle := subtitleStyle.Render("Configure sound alerts")

	items := m.soundMenuItems()
	lines := make([]string, len(items))
	for i, item := range items {
		lines[i] = m.formatSoundLine(i, item)
	}

	help := helpStyle.Render("[j/k] navigate  [enter] select  [esc] back  [q] quit")

	content := strings.Join(lines, "\n")

	return containerStyle.Render(
		lipgloss.JoinVertical(lipgloss.Center, title, subtitle, "", content, "", help),
	)
}

func (m Model) formatSoundLine(idx int, text string) string {
	prefix := "  "
	style := itemStyle
	if m.cursor == idx {
		prefix = "> "
		style = selectedItemStyle
	}
	return style.Render(prefix + text)
}

func (m Model) soundMenuItems() []string {
	enabledValue := "Off"
	if m.cfg.Sound.Enabled {
		enabledValue = "On"
	}

	modeLabel := "Terminal bell"
	if m.cfg.Sound.Mode == config.SoundModeMac {
		modeLabel = "macOS system sound"
	}

	items := []string{
		fmt.Sprintf("Enabled: %s", enabledValue),
		fmt.Sprintf("Mode: %s", modeLabel),
	}

	if m.cfg.Sound.Mode == config.SoundModeMac {
		tone := m.cfg.Sound.Tone
		if tone == "" {
			tone = config.DefaultMacTone
		}
		items = append(items, fmt.Sprintf("Tone: %s", tone))
	}

	items = append(items, "Test sound", "Back")
	return items
}

func macTones() []string {
	return []string{
		"Basso",
		"Blow",
		"Bottle",
		"Frog",
		"Funk",
		"Glass",
		"Hero",
		"Morse",
		"Ping",
		"Pop",
		"Purr",
		"Sosumi",
		"Submarine",
		"Tink",
	}
}
