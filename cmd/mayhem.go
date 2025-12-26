// Package main is the core application
package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
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
	forceImport := false
	verbose := false
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
		isVerbose := set.Bool("verbose", false, "enable verbose output")
		var force *bool
		if isImport {
			force = set.Bool("overwrite", false, "force import to overwrite current data")
		}
		if err := set.Parse(args); err != nil {
			return err
		}
		if force != nil {
			forceImport = *force
		}
		if isVerbose != nil {
			verbose = *isVerbose
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
		if !forceImport {
			return errors.New("import not supported into existing database")
		}
		if err := os.Remove(ctx.Config.Database()); err != nil {
			return err
		}
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

	errors := ctx.DB.Errors()
	if len(errors) > 0 {
		f, err := os.OpenFile(filepath.Join(ctx.Config.Data.Directory, "log.txt"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
		if err != nil {
			return err
		}
		defer f.Close()
		for _, item := range errors {
			if verbose {
				fmt.Fprintf(os.Stderr, "%s\n", item)
			}
			fmt.Fprintf(f, "%s\n", item)
		}
	}
	return err
}
