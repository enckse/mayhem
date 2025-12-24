// Package state handles overall app state
package state

import (
	"errors"
	"os"

	"github.com/enckse/mayhem/internal/display"
	"gorm.io/gorm"
)

// FileName is the application name prefix
const FileName = "todo."

// Context is the overall state context
type Context struct {
	DB     *gorm.DB
	Config Config
	Screen *display.Screen
	State  struct {
		FinishedTasks map[uint]bool
	}
}

// PathExists will indicate if a path exists (or not)
func PathExists(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}
