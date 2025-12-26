package entities

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/enckse/mayhem/internal/backend"
	"github.com/google/uuid"
)

const maxPriority = 4

// MaxPriority is the maximum allowed priority
var MaxPriority = fmt.Sprintf("%d", maxPriority)

// Task defines task-based entities for work
type Task struct {
	ID       string
	Title    string
	Notes    string
	Deadline time.Time
	Priority uint64
	Finished time.Time
	StackID  string
}

// NewTask will create a new task
func NewTask() Task {
	return Task{ID: uuid.NewString()}
}

// Save will store the task
func (t Task) Save(store backend.Store) Entity {
	if strings.TrimSpace(t.Title) == "" {
		store.Log("task", errors.New("no title"))
		return t
	}
	if t.Priority > maxPriority {
		store.Log("task", errors.New("invalid priority"))
		return t
	}
	store.AddChild(t.StackID, t.ID, t)
	return t
}

// Delete will remove the task
func (t Task) Delete(store backend.Store) {
	store.RemoveChild(t.StackID, t.ID)
}

// SortTasks will sort by finished, deadline, title
func SortTasks(t []Task) {
	slices.SortFunc(t, func(x, y Task) int {
		if !x.Finished.IsZero() && y.Finished.IsZero() {
			return 1
		}
		if x.Finished.IsZero() && !y.Finished.IsZero() {
			return -1
		}
		if val := x.Finished.Compare(y.Finished); val != 0 {
			return val
		}
		if x.Deadline.IsZero() && !y.Deadline.IsZero() {
			return 1
		}
		if !x.Deadline.IsZero() && y.Deadline.IsZero() {
			return -1
		}
		if val := x.Deadline.Compare(y.Deadline); val != 0 {
			return val
		}
		return strings.Compare(x.Title, y.Title)
	})
}

// EntityID will get the entity ID for the object
func (t Task) EntityID() string {
	return t.ID
}
