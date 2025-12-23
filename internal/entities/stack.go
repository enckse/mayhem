package entities

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Stack is a set of tasks, sorted alphabetically
type Stack struct {
	gorm.Model
	Title            string `gorm:"notnull"`
	PendingTaskCount int
	Tasks            []Task
}

// EntityID gets the backing entity id
func (s Stack) EntityID() uint {
	return s.ID
}

// InitializeStacks will initialize the initial stack set
func InitializeStacks() (Stack, error) {
	stack := Stack{Title: "New Stack"}
	result := DB.Create(&stack)
	return stack, result.Error
}

// FetchAllStacks will retrieve all stacks
func FetchAllStacks() ([]Stack, error) {
	var stacks []Stack
	result := DB.Model(&Stack{}).Preload("Tasks").Find(&stacks)

	if len(stacks) == 0 {
		stack, err := InitializeStacks()
		return []Stack{stack}, err
	}

	return stacks, result.Error
}

// IncPendingCount will add to the pending task count
func IncPendingCount(id uint) {
	stack := Stack{}
	DB.Find(&stack, id)
	stack.PendingTaskCount++
	stack.Save()
}

// Save will save the entity
func (s Stack) Save() Entity {
	DB.Save(&s)
	return s
}

// Delete will remove the entity
func (s Stack) Delete() {
	// Unscoped() is used to ensure hard delete, where stack will be removed from db instead of being just marked as "deleted"
	// DB.Unscoped().Delete(&s)
	DB.Unscoped().Select(clause.Associations).Delete(&s)
}
