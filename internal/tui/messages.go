// Package tui requires TUI messages
package tui

import tea "github.com/charmbracelet/bubbletea"

type goToMainMsg struct {
	value any
}

func goToMainCmd() tea.Msg {
	return goToMainMsg{
		value: "",
	}
}

func goToMainWithVal(value any) tea.Cmd {
	return func() tea.Msg {
		return goToMainMsg{value: value}
	}
}

type goToFormMsg struct {
	value any
}

func goToFormWithVal(value any) tea.Cmd {
	return func() tea.Msg {
		return goToFormMsg{value: value}
	}
}

type goToStepsMsg struct {
	value any
}

func goToStepsWithVal(value any) tea.Cmd {
	return func() tea.Msg {
		return goToStepsMsg{value: value}
	}
}
