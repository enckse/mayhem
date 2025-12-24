// Package main is the core application
package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/enckse/mayhem/internal/app"
	"github.com/enckse/mayhem/internal/convert"
	entities "github.com/enckse/mayhem/internal/entities"
	"github.com/enckse/mayhem/internal/state"
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
	ctx := &state.Context{}
	if err := entities.InitializeDB(ctx); err != nil {
		return err
	}

	model := tui.InitializeMainModel(ctx)
	p := tea.NewProgram(model.Backing, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		return err
	}
	if strings.TrimSpace(os.Getenv(app.EnvPrefix+"EXPORT_JSON")) == "1" {
		return convert.ToJSON(ctx)
	}
	return nil
}
