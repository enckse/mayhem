// Package state handles overall app state
package state

import (
	"errors"
	"os"

	"github.com/enckse/mayhem/internal/display"
)

// FileName is the application name prefix
const FileName = "todo."

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

// PathExists will indicate if a path exists (or not)
func PathExists(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}
