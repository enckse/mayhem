// Package textarea is a large text area input
package textarea

import (
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/enckse/mayhem/internal/display"
	"github.com/enckse/mayhem/internal/tui/messages"
)

// Input is a textarea input field
// textarea.Model doesn't implement tea.Model interface
type Input struct {
	input textarea.Model
}

// New will create a new textarea
func New(value string, screen *display.Screen) tea.Model {
	t := textarea.New()
	t.SetValue(value)
	t.SetWidth(screen.InputFormStyle().GetWidth() - 2)
	t.SetHeight(4)
	t.CharLimit = 500
	t.Placeholder = "Enter task notes"
	t.ShowLineNumbers = false
	// We only deal with textarea in focused state, so blurred style is redundant
	t.FocusedStyle = textarea.Style{Placeholder: display.PlaceHolderStyle, Text: display.TextInputStyle}
	t.Focus()

	m := Input{
		input: t,
	}

	return m
}

// Init will init the model
func (m Input) Init() tea.Cmd {
	return textarea.Blink
}

// Update will update the model
func (m Input) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+s":
			return m, messages.FormGoToWith(m.input.Value())
		}
	}

	// Placing it outside KeyMsg case is required, otherwise messages like textarea's Blink will be lost
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

// View will display the view
func (m Input) View() string {
	// Can't just render textarea.Value(), otherwise cursor blinking wouldn't work
	return m.input.View()
}
