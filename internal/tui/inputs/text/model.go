// Package text defines simple text inputs
package text

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/enckse/mayhem/internal/display"
	"github.com/enckse/mayhem/internal/tui/keys"
)

// Input defines the text input type
// textinput.Model doesn't implement tea.Model interface
type Input struct {
	input textinput.Model
	// Since textinput field can be used in multiple places,
	// responder is required to determine the receiver of the message emitted by textinput field
	responder func(any) tea.Cmd
}

// New creates a new text entry input
func New(value, placeholder string, charLimit int, responder func(any) tea.Cmd) tea.Model {
	t := textinput.New()
	t.SetValue(value)

	t.Cursor.Style = display.TextInputStyle
	t.CharLimit = charLimit
	t.Focus()
	t.PromptStyle = display.TextInputStyle
	t.TextStyle = display.TextInputStyle
	t.Placeholder = placeholder

	m := Input{
		input:     t,
		responder: responder,
	}

	return m
}

// Init will init the model
func (m Input) Init() tea.Cmd {
	return textinput.Blink
}

// Update will update the model
func (m Input) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Mappings.Enter):
			return m, m.responder(m.input.Value())
		}
	}

	// Placing it outside KeyMsg case is required, otherwise messages like textinput's Blink will be lost
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

// View will view the model
func (m Input) View() string {
	// Can't just render textinput.Value(), otherwise cursor blinking wouldn't work
	return m.input.View()
}
