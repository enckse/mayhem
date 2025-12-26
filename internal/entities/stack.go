package entities

import (
	"errors"
	"slices"
	"strings"

	"github.com/enckse/mayhem/internal/backend"
	"github.com/google/uuid"
)

// Stack is a set of tasks, sorted alphabetically
type Stack struct {
	ID    string
	Title string
	Tasks []Task `json:"-"`
}

// OpenTasks will get the count of unfinished tasks
func (s Stack) OpenTasks() uint64 {
	var count uint64
	for _, t := range s.Tasks {
		if !t.Finished.IsZero() {
			continue
		}
		count++
	}
	return count
}

// NewStack will create a new stack
func NewStack(store backend.Store) Stack {
	stack := Stack{Title: "New Stack"}
	stack.ID = uuid.NewString()
	store.Add(stack.ID, stack)
	return stack
}

// FetchStacks will retrieve all stacks
func FetchStacks(store backend.Store) []Stack {
	var stacks []Stack
	for _, item := range store.Get() {
		c, ok := item.Node.(Stack)
		if !ok {
			continue
		}
		c.Tasks = []Task{}
		for _, t := range item.Children {
			task, ok := t.Node.(Task)
			if ok {
				c.Tasks = append(c.Tasks, task)
			}
		}
		stacks = append(stacks, c)
	}

	if len(stacks) == 0 {
		stack := NewStack(store)
		return []Stack{stack}
	}

	return stacks
}

// EntityID will get the entity ID for the object
func (s Stack) EntityID() string {
	return s.ID
}

// Save will save the entity
func (s Stack) Save(store backend.Store) Entity {
	if strings.TrimSpace(s.Title) == "" {
		store.Log("stack", errors.New("no title"))
		return s
	}
	store.Add(s.ID, s)
	return s
}

// Delete will remove the entity
func (s Stack) Delete(store backend.Store) {
	store.Remove(s.ID)
}

// SortStacks will sort by title
func SortStacks(s []Stack) {
	slices.SortFunc(s, func(x, y Stack) int {
		return strings.Compare(x.Title, y.Title)
	})
}
