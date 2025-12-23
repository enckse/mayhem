package entities

import (
	"time"

	"gorm.io/gorm"
)

// RecurTask is a recuring task
type RecurTask struct {
	gorm.Model
	Deadline   time.Time `gorm:"index:idx_member"`
	IsFinished bool
	StackID    uint `gorm:"index:idx_member"`
	TaskID     uint
}

// Save will save the task
func (r RecurTask) Save() {
	DB.Save(&r)
}
