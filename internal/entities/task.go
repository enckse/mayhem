package entities

import (
	"slices"
	"strings"
	"time"

	"github.com/enckse/mayhem/internal/state"
	"gorm.io/gorm"
)

// Task defines task-based entities for work
type Task struct {
	gorm.Model `json:"-"`
	Title      string `gorm:"notnull"`
	Notes      string
	Deadline   time.Time
	Priority   uint64
	IsFinished bool
	StackID    uint `json:"-"`
}

// Save will store the task
func (t Task) Save(store state.Store) Entity {
	store.Save(&t)
	return t
}

// Delete will remove the task
func (t Task) Delete(store state.Store) {
	store.Delete(&t)
}

// EntityID gets the backing entity id
func (t Task) EntityID() uint {
	return t.ID
}

// SortTasks will sort by finished, deadline, title
func SortTasks(t []Task) {
	slices.SortFunc(t, func(x, y Task) int {
		if x.IsFinished && !y.IsFinished {
			return -1
		}
		if !x.IsFinished && y.IsFinished {
			return 1
		}
		if val := x.Deadline.Compare(y.Deadline); val != 0 {
			return val
		}
		return strings.Compare(x.Title, y.Title)
	})
}
