// Package deletion handles confirmation of removal
package deletion

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/enckse/mayhem/internal/display"
	"github.com/enckse/mayhem/internal/tui/definitions"
	"github.com/enckse/mayhem/internal/tui/keys"
	"github.com/enckse/mayhem/internal/tui/messages"
)

const isConfirm = "y"

// Confirmation is a delete confirmation handler
// textinput.Model doesn't implement tea.Model interface
type Confirmation struct {
	customInputType string
}

// NewConfirmation creates a new confirmation model item
func NewConfirmation() tea.Model {
	m := Confirmation{
		customInputType: definitions.IsDelete,
	}

	return m
}

// Init will initialize the object
func (m Confirmation) Init() tea.Cmd {
	return textinput.Blink
}

// Update will update the object
func (m Confirmation) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Mappings.Return):
			return m, messages.MainGoTo
		case key.Matches(msg, keys.Mappings.Quit):
			return m, tea.Quit
		default:
			if strings.ToLower(msg.String()) == isConfirm {
				return m, messages.MainGoToWith(isConfirm)
			}
			return m, messages.MainGoToWith("")
		}
	}
	return m, nil
}

// View will handle rendering the view
func (m Confirmation) View() string {
	// Can't just render textinput.Value(), otherwise cursor blinking wouldn't work
	return lipgloss.NewStyle().Foreground(display.HighlightedBackgroundColor).Padding(1, 0).Render("Do you wish to proceed with deletion? (" + isConfirm + "/n): ")
}
