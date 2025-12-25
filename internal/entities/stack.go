package entities

import (
	"slices"
	"strings"

	"github.com/enckse/mayhem/internal/state"
	"gorm.io/gorm"
)

// Stack is a set of tasks, sorted alphabetically
type Stack struct {
	gorm.Model `json:"-"`
	Title      string `gorm:"notnull"`
	Tasks      []Task
}

// OpenTasks will get the count of unfinished tasks
func (s Stack) OpenTasks() uint64 {
	var count uint64
	for _, t := range s.Tasks {
		if t.IsFinished {
			continue
		}
		count++
	}
	return count
}

// EntityID gets the backing entity id
func (s Stack) EntityID() uint {
	return s.ID
}

// NewStack will create a new stack
func NewStack(store state.Store) Stack {
	stack := Stack{Title: "New Stack"}
	store.Create(&stack)
	return stack
}

// FetchStacks will retrieve all stacks
func FetchStacks(store state.Store) []Stack {
	obj := store.Fetch()
	var stacks []Stack
	val, ok := obj.([]Stack)
	if ok {
		stacks = val
	}

	if len(stacks) == 0 {
		stack := NewStack(store)
		return []Stack{stack}
	}

	return stacks
}

// Save will save the entity
func (s Stack) Save(store state.Store) Entity {
	store.Save(&s)
	return s
}

// Delete will remove the entity
func (s Stack) Delete(store state.Store) {
	store.Delete(&s)
}

// SortStacks will sort by title
func SortStacks(s []Stack) {
	slices.SortFunc(s, func(x, y Stack) int {
		return strings.Compare(x.Title, y.Title)
	})
}
