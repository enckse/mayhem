// Package main is the core application
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/enckse/mayhem/internal/backend"
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
	ctx := &state.Context{}
	ctx.Screen = display.NewScreen()
	cfg, err := state.LoadConfig(configFile)
	if err != nil {
		return err
	}
	ctx.Config = cfg
	if ctx.Config.Backups.Directory != "" {
		if err := ctx.Config.Backup(time.Now()); err != nil {
			return err
		}
	}
	file := ctx.Config.Database()
	f, err := os.OpenFile(filepath.Join(ctx.Config.Data.Directory, "log.txt"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	storage := backend.NewMemoryBased(file, ctx.Config.Data.Pretty, f)
	if state.PathExists(file) {
		if err := backend.Load[entities.Stack, entities.Task](storage); err != nil {
			return err
		}
	}
	ctx.DB = storage
	err = func() error {
		model := ui.Initialize(ctx)
		p := tea.NewProgram(model.Backing, tea.WithAltScreen())

		if _, err := p.Run(); err != nil {
			return err
		}
		return nil
	}()

	return err
}
