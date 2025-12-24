package entities

import (
	"slices"
	"strings"

	"github.com/enckse/mayhem/internal/state"
	"gorm.io/gorm"
)

// Stack is a set of tasks, sorted alphabetically
type Stack struct {
	gorm.Model       `json:"-"`
	Title            string `gorm:"notnull"`
	PendingTaskCount uint64 `json:"-"`
	Tasks            []Task
}

// EntityID gets the backing entity id
func (s Stack) EntityID() uint {
	return s.ID
}

// NewStack will create a new stack
func NewStack(ctx *state.Context) (Stack, error) {
	stack := Stack{Title: "New Stack"}
	err := ctx.DB.Create(&stack)
	return stack, err
}

// FetchStacks will retrieve all stacks
func FetchStacks(ctx *state.Context) ([]Stack, error) {
	var stacks []Stack
	err := ctx.DB.Stacks(&stacks)
	if err != nil {
		return stacks, err
	}

	if len(stacks) == 0 {
		stack, err := NewStack(ctx)
		return []Stack{stack}, err
	}

	return stacks, err
}

// IncrementPendingCount will add to the pending task count
func IncrementPendingCount(id uint, ctx *state.Context) {
	stack := Stack{}
	ctx.DB.Find(&stack, id)
	stack.PendingTaskCount++
	stack.Save(ctx)
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
