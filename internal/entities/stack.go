package entities

import (
	"slices"
	"strings"

	"github.com/enckse/mayhem/internal/state"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
	result := ctx.DB.Create(&stack)
	return stack, result.Error
}

// FetchStacks will retrieve all stacks
func FetchStacks(ctx *state.Context) ([]Stack, error) {
	var stacks []Stack
	result := ctx.DB.Model(&Stack{}).Preload("Tasks").Find(&stacks)

	if len(stacks) == 0 {
		stack, err := NewStack(ctx)
		return []Stack{stack}, err
	}

	return stacks, result.Error
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
	// Unscoped() is used to ensure hard delete, where stack will be removed from db instead of being just marked as "deleted"
	// DB.Unscoped().Delete(&s)
	ctx.DB.Unscoped().Select(clause.Associations).Delete(&s)
}

// SortStacks will sort by title
func SortStacks(s []Stack) {
	slices.SortFunc(s, func(x, y Stack) int {
		return strings.Compare(x.Title, y.Title)
	})
}
