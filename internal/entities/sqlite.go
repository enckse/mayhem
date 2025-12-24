package entities

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DBWrapper struct{ db *gorm.DB }

func (d *DBWrapper) Save(obj any) {
	d.db.Save(obj)
}

func (d *DBWrapper) Create(obj any) error {
	return d.db.Create(obj).Error
}

func (d *DBWrapper) Stacks(obj any) error {
	result := d.db.Model(&Stack{}).Preload("Tasks").Find(obj)
	return result.Error
}

func (d *DBWrapper) Find(obj any, id uint) {
	d.db.Find(obj, id)
}

func (d *DBWrapper) Delete(obj any) {
	// Unscoped() is used to ensure hard delete, where stack will be removed from db instead of being just marked as "deleted"
	// DB.Unscoped().Delete(&s)
	d.db.Unscoped().Select(clause.Associations).Delete(obj)
}
