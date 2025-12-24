// Package main is the core application
package main

import (
	"errors"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/enckse/mayhem/internal/display"
	"github.com/enckse/mayhem/internal/entities"
	"github.com/enckse/mayhem/internal/state"
	"github.com/enckse/mayhem/internal/tui"
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
	isMerge := false
	for idx, arg := range os.Args {
		switch idx {
		case 0:
			continue
		case 1:
		default:
			return fmt.Errorf("unexpected argument: %s", arg)
		}
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
		case "merge":
			isImport = true
			isMerge = true
			continue
		}
	}
	if isExport && isImport {
		return errors.New("only one of export/import can be provided")
	}
	ctx := &state.Context{}
	ctx.Screen = display.NewScreen()
	cfg, err := state.LoadConfig()
	if err != nil {
		return err
	}
	ctx.Config = cfg
	exists := state.PathExists(ctx.Config.Database())
	if isExport && !exists {
		return errors.New("no database to dump")
	}
	if isImport && exists && !isMerge {
		return errors.New("import not supported into existing database")
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
		return entities.DumpJSON(ctx)
	}
	if isImport {
		return entities.LoadJSON(ctx, isMerge)
	}

	model := tui.InitializeMainModel(ctx)
	p := tea.NewProgram(model.Backing, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		return err
	}
	if ctx.Config.JSON.Exit {
		return entities.ToJSON(ctx)
	}
	return nil
}
