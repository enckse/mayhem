package entities

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Task defines task-based entities for work
type Task struct {
	gorm.Model
	Title              string `gorm:"notnull"`
	Description        string
	Deadline           time.Time
	Priority           int // 3: High, 2: Mid, 1: Low, 0: No Priority
	IsFinished         bool
	StackID            uint
}

// Save will store the task
func (t Task) Save() Entity {
	DB.Save(&t)
	return t
}

// Delete will remove the task
func (t Task) Delete() {
	// Unscoped() is used to ensure hard delete, where task will be removed from db instead of being just marked as "deleted"
	DB.Unscoped().Select(clause.Associations).Delete(&t)
}

// EntityID gets the backing entity id
func (t Task) EntityID() uint {
	return t.ID
}
