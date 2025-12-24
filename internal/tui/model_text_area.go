package tui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/enckse/mayhem/internal/display"
	"github.com/enckse/mayhem/internal/tui/keys"
)

// textarea.Model doesn't implement tea.Model interface
type textArea struct {
	input textarea.Model
}

var textAreaKeys = keys.Map{
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("'enter'", "new line"),
	),
	Save: key.NewBinding(
		key.WithKeys("ctrl+s"),
		key.WithHelp("'ctrl+s'", "save"),
	),
	Return: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("'esc'", "return"),
	),
}

func initializeTextArea(value string, screen *display.Screen) tea.Model {
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

	m := textArea{
		input: t,
	}

	return m
}

func (m textArea) Init() tea.Cmd {
	return textarea.Blink
}

func (m textArea) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+s":
			return m, goToFormWithVal(m.input.Value())
		}
	}

	// Placing it outside KeyMsg case is required, otherwise messages like textarea's Blink will be lost
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m textArea) View() string {
	// Can't just render textarea.Value(), otherwise cursor blinking wouldn't work
	return m.input.View()
}
