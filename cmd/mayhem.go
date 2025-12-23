package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	entities "github.com/enckse/mayhem/internal/entities"
	tui "github.com/enckse/mayhem/internal/tui"
)

var version string

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func run() error {
	for idx, arg := range os.Args {
		if idx == 0 {
			continue
		}
		if arg == "--version" {
			fmt.Fprintf(os.Stderr, "%s\n", version)
			return nil
		}
		return fmt.Errorf("invalid arguments: %s", arg)
	}
	if err := entities.InitializeDB(); err != nil {
		return err
	}

	model := tui.InitializeMainModel()
	p := tea.NewProgram(model, tea.WithAltScreen())

	_, err := p.Run()
	return err
}
