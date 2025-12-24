// Package lists has list selectors
package lists

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/enckse/mayhem/internal/display"
	"github.com/enckse/mayhem/internal/tui/definitions"
	"github.com/enckse/mayhem/internal/tui/keys"
	"github.com/enckse/mayhem/internal/tui/messages"
)

// Selector is a list picker/selector
type Selector struct {
	options    []definitions.KeyValue
	focusIndex int
	maxIndex   int
	responder  func(any) tea.Cmd
}

// Init will init the model
func (m Selector) Init() tea.Cmd {
	return nil
}

// NewSelector will create a new list selector
func NewSelector(options []definitions.KeyValue, selectedVal string, responder func(any) tea.Cmd) tea.Model {
	// Takes care of default case where index should be 0
	var selectedIndex int

	for i, item := range options {
		if item.Value == selectedVal {
			selectedIndex = i
			break
		}
	}

	m := Selector{
		focusIndex: selectedIndex,
		maxIndex:   len(options) - 1,
		options:    options,
		responder:  responder,
	}

	return m
}

// Update will update the model
func (m Selector) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Mappings.Return):
			return m, messages.MainGoToWith(definitions.KeyValue{})

		case key.Matches(msg, keys.Mappings.Quit, keys.Mappings.Exit):
			return m, tea.Quit

		case key.Matches(msg, keys.Mappings.Enter):
			return m, m.responder(m.options[m.focusIndex])

		case key.Matches(msg, keys.Mappings.Up):
			if m.focusIndex > 0 {
				m.focusIndex--
			} else {
				m.focusIndex = m.maxIndex
				return m, nil
			}

		case key.Matches(msg, keys.Mappings.Down):
			if m.focusIndex < m.maxIndex {
				m.focusIndex++
			} else {
				m.focusIndex = 0
				return m, nil
			}
		}
	}
	return m, nil
}

// View will show the model
func (m Selector) View() string {
	var res []string

	for i, item := range m.options {
		var value string

		if i == m.focusIndex {
			value = lipgloss.NewStyle().Foreground(display.InputFormColor).Bold(true).Render("» " + item.Value)
		} else {
			value = lipgloss.NewStyle().Foreground(display.InputFormColor).Bold(true).Render("  " + item.Value)
		}

		res = append(res, value)
	}

	return lipgloss.JoinVertical(lipgloss.Left, res...)
}
