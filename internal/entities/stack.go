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
func NewStack(ctx *state.Context) Stack {
	stack := Stack{Title: "New Stack"}
	ctx.DB.Create(&stack)
	return stack
}

// FetchStacks will retrieve all stacks
func FetchStacks(ctx *state.Context) []Stack {
	obj := ctx.DB.Fetch()
	var stacks []Stack
	val, ok := obj.([]Stack)
	if ok {
		stacks = val
	}

	if len(stacks) == 0 {
		stack := NewStack(ctx)
		return []Stack{stack}
	}

	return stacks
}

// Save will save the entity
func (s Stack) Save(ctx *state.Context) Entity {
	Sync(ctx, &s)
	return s
}

// Delete will remove the entity
func (s Stack) Delete(ctx *state.Context) {
	ctx.DB.Delete(&s)
}

// SortStacks will sort by title
func SortStacks(s []Stack) {
	slices.SortFunc(s, func(x, y Stack) int {
		return strings.Compare(x.Title, y.Title)
	})
}
