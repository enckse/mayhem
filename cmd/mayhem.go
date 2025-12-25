// Package main is the core application
package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/enckse/mayhem/internal/display"
	"github.com/enckse/mayhem/internal/entities"
	"github.com/enckse/mayhem/internal/state"
	"github.com/enckse/mayhem/internal/tui/ui"
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
	args := os.Args
	var configFile string
	if len(args) > 1 {
		args = args[1:]
		cmd := args[0]
		if !strings.HasPrefix(cmd, "-") {
			switch cmd {
			case "version":
				fmt.Fprintf(os.Stderr, "%s\n", version)
				return nil
			case "export":
				isExport = true
			case "import":
				isImport = true
			case "merge":
				isImport = true
				isMerge = true
			}
			if len(args) > 1 {
				args = args[1:]
			}
		}
		set := flag.NewFlagSet("cli", flag.ExitOnError)
		cfgFile := set.String("config", "", "configuration file")
		if err := set.Parse(args); err != nil {
			return err
		}
		configFile = *cfgFile
	}
	if isExport && isImport {
		return errors.New("only one of export/import can be provided")
	}
	ctx := &state.Context{}
	ctx.Screen = display.NewScreen()
	cfg, err := state.LoadConfig(configFile)
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
		if ctx.Config.Backups.Directory != "" {
			if err := ctx.Config.Backup(time.Now()); err != nil {
				return err
			}
		}
	}
	if err := entities.InitializeDB(ctx); err != nil {
		return err
	}
	err = func() error {
		if isExport {
			return entities.DumpJSON(os.Stdout, ctx.DB)
		}
		if isImport {
			return entities.LoadJSON(ctx.DB, isMerge, os.Stdin)
		}

		model := ui.Initialize(ctx)
		p := tea.NewProgram(model.Backing, tea.WithAltScreen())

		if _, err := p.Run(); err != nil {
			return err
		}
		return nil
	}()
	for _, item := range ctx.DB.Errors() {
		fmt.Fprintf(os.Stderr, "%s\n", item)
	}
	return err
}
