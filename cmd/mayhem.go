// Package main is the core application
package main

import (
	"errors"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
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
	isExport := false
	isImport := false
	for idx, arg := range os.Args {
		if idx == 0 {
			continue
		}
		switch arg {
		case "--version":
			fmt.Fprintf(os.Stderr, "%s\n", version)
			return nil
		case "export":
			isExport = true
			continue
		case "import":
			isImport = true
			continue
		}

		return fmt.Errorf("invalid arguments: %s", arg)
	}
	if isExport && isImport {
		return errors.New("only one of export/import can be provided")
	}
	ctx := &state.Context{}
	cfg, err := state.LoadConfig()
	if err != nil {
		return err
	}
	ctx.Config = cfg
	if isExport && !state.PathExists(ctx.Config.Database()) {
		return errors.New("no database to dump")
	}
	if !isExport && !isImport {
		if err := ctx.Config.Backup(); err != nil {
			return err
		}
	}
	if err := entities.InitializeDB(ctx); err != nil {
		return err
	}
	if isExport {
		return convert.DumpJSON(ctx)
	}

	model := tui.InitializeMainModel(ctx)
	p := tea.NewProgram(model.Backing, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		return err
	}
	if ctx.Config.Convert.JSON {
		return convert.ToJSON(ctx)
	}
	return nil
}
