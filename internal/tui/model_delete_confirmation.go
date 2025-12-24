package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/enckse/mayhem/internal/tui/display"
	"github.com/enckse/mayhem/internal/tui/keys"
)

const isConfirm = "y"

// textinput.Model doesn't implement tea.Model interface
type deleteConfirmation struct {
	customInputType string
}

func initializeDeleteConfirmation() tea.Model {
	m := deleteConfirmation{
		customInputType: "delete",
	}

	return m
}

func (m deleteConfirmation) Init() tea.Cmd {
	return textinput.Blink
}

func (m deleteConfirmation) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {

		case key.Matches(msg, keys.Mappings.Return):
			return m, goToMainCmd

		case key.Matches(msg, keys.Mappings.Quit):
			return m, tea.Quit

		default:
			if strings.ToLower(msg.String()) == isConfirm {
				return m, goToMainWithVal(isConfirm)
			}
			return m, goToMainWithVal("")
		}
	}
	return m, nil
}

func (m deleteConfirmation) View() string {
	// Can't just render textinput.Value(), otherwise cursor blinking wouldn't work
	return lipgloss.NewStyle().Foreground(display.HighlightedBackgroundColor).Padding(1, 0).Render("Do you wish to proceed with deletion? (" + isConfirm + "/n): ")
}
