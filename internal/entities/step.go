package entities

import (
	"gorm.io/gorm"
)

// Step is a task step
type Step struct {
	gorm.Model
	Title      string
	IsFinished bool
	TaskID     uint
}

// Save will store the step
func (s Step) Save() Step {
	DB.Save(&s)
	return s
}

// Delete will remove the step
func (s Step) Delete() {
	DB.Delete(&s)
}
