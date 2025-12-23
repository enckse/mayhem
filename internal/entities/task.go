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
	Steps              []Step
	Deadline           time.Time
	Priority           int // 3: High, 2: Mid, 1: Low, 0: No Priority
	IsFinished         bool
	IsRecurring        bool
	StartTime          time.Time // Applicable only for recurring tasks
	RecurrenceInterval int       // in days
	RecurChildren      []RecurTask
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

// LatestRecurTask will get the latest recurring task
func (t Task) LatestRecurTask() (RecurTask, int64) {
	recurTask := RecurTask{}
	// localtime modifier has to be added to DATE other wise UTC time would be used
	result := DB.Last(&recurTask, "task_id = ? AND deadline <  DATE('now', 'localtime', 'start of day', '+1 day')", t.ID)
	return recurTask, result.RowsAffected
}

// RemoveFutureRecurTasks will remove future recurring tasks
func (t Task) RemoveFutureRecurTasks() {
	DB.Unscoped().Where("deadline >=  DATE('now', 'start of day') AND task_id = ?", t.ID).Delete(&RecurTask{})
}

// FetchAllRecurTasks will fetch all recurring tasks
func (t Task) FetchAllRecurTasks() []RecurTask {
	DB.Preload("RecurChildren").Find(&t)
	return t.RecurChildren
}

func (t Task) EntityID() uint {
	return t.ID
}
