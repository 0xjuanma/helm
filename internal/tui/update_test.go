package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/0xjuanma/helm/internal/config"
)

func newSelectModelAtCursor(cursor int) Model {
	cfg := &config.Config{}
	return Model{
		screen:    screenSelect,
		workflows: cfg.BuildWorkflows(),
		cfg:       cfg,
		cursor:    cursor,
	}
}

func TestHandleSelectKeyNoOpOnEmptySlots(t *testing.T) {
	for _, cursor := range []int{2, 3, 4} {
		m := newSelectModelAtCursor(cursor)

		got, _ := m.handleSelectKey(tea.KeyMsg{Type: tea.KeyEnter})

		gotModel := got.(Model)
		if gotModel.screen != screenSelect {
			t.Errorf("cursor %d: screen = %v, want screenSelect (enter on empty slot should be a no-op)", cursor, gotModel.screen)
		}
	}
}
