package entities

import (
	"github.com/enckse/mayhem/internal/state"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// DBWrapper is our sqlite backend wrapper
type DBWrapper struct{ db *gorm.DB }

// Save will store an entity
func (d *DBWrapper) Save(obj any) {
	d.db.Save(obj)
}

// Create will create an entity
func (d *DBWrapper) Create(obj any) error {
	return d.db.Create(obj).Error
}

// Stacks will get the stacks and preload tasks
func (d *DBWrapper) Stacks(obj any) error {
	result := d.db.Model(&Stack{}).Preload("Tasks").Find(obj)
	return result.Error
}

// Find will find an entity in the database
func (d *DBWrapper) Find(obj any, id uint) {
	d.db.Find(obj, id)
}

// Delete will remove an object
// Unscoped() is used to ensure hard delete, where stack will be removed from db instead of being just marked as "deleted"
// DB.Unscoped().Delete(&s)
func (d *DBWrapper) Delete(obj any) {
	d.db.Unscoped().Select(clause.Associations).Delete(obj)
}

// SyncJSON will sync db state to JSON
func (d *DBWrapper) SyncJSON(ctx *state.Context) {
	ToJSON(ctx)
}
