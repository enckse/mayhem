// Package help implements the help view models
package help

import (
	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/enckse/mayhem/internal/tui/keys"
)

// Model is the underlying help model
type Model struct {
	help help.Model
	keys keys.Map
}

// NewModel will create a new help model
func NewModel(keys keys.Map) Model {
	return Model{
		keys: keys,
		help: help.New(),
	}
}

// Init is the model init
func (m Model) Init() tea.Cmd {
	return nil
}

// Update will update the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// If we set a width on the help menu it can it can gracefully truncate
		// its view as needed.
		m.help.Width = msg.Width
	}

	return m, nil
}

// View will handle view rendering
func (m Model) View() string {
	style := lipgloss.NewStyle().MarginTop(1)
	return style.Render(m.help.View(m.keys))
}
