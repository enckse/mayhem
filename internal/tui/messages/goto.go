// Package messages allows form moving between models/UI elements
package messages

import tea "github.com/charmbracelet/bubbletea"

type (
	// Main will move to main view
	Main struct {
		Value any
	}
	// Form will move to a form
	Form struct {
		Value any
	}
)

// MainGoTo will go to the main via empty message
func MainGoTo() tea.Msg {
	return Main{
		Value: "",
	}
}

// MainGoToWith will go to main with a value
func MainGoToWith(value any) tea.Cmd {
	return func() tea.Msg {
		return Main{Value: value}
	}
}

// FormGoToWith will go to a form with a value
func FormGoToWith(value any) tea.Cmd {
	return func() tea.Msg {
		return Form{Value: value}
	}
}
