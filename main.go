package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	entities "github.com/enckse/mayhem/entities"
	tui "github.com/enckse/mayhem/tui"
)

func main() {
	entities.InitializeDB()

	model := tui.InitializeMainModel()
	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		log.Fatal("Error encountered while running the program:", err)
	}
}
