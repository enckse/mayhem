// Package state handles overall app state
package state

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/enckse/mayhem/internal/display"
)

const (
	// FileName is the application name prefix
	FileName  = "todo."
	envPrefix = "MAYHEM_"
)

type (
	// Store defines the backing store for data
	Store interface {
		Save(any)
		Create(any) error
		Stacks(any) error
		Find(any, uint)
		Delete(any)
		SyncJSON(*Context)
	}
	// Context is the overall state context
	Context struct {
		DB     Store
		Config Config
		Screen *display.Screen
		State  struct {
			FinishedTasks map[uint]bool
		}
	}
)

// PathExists indicates if a path exists
func PathExists(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

func detectDir(xdgName, envVar, altName string) (string, error) {
	p, err := getDir(xdgName, envVar, altName)
	if err != nil {
		return "", err
	}
	return p, os.MkdirAll(p, os.ModePerm)
}

func getDir(xdgName, envVar, altName string) (string, error) {
	path := os.Getenv(envPrefix + envVar)
	if path != "" {
		return path, nil
	}

	const appDir = "mayhem"
	xdg := os.Getenv(xdgName)
	if xdg != "" {
		return filepath.Join(xdg, appDir), nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, altName, appDir), nil
}
