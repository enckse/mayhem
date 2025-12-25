package entities

import (
	"github.com/enckse/mayhem/internal/state"
	// Using pure-go implementation of GORM driver to avoid CGO issues during cross-compilation
	"github.com/enckse/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
	_ "modernc.org/sqlite"
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

// Fetch will return the data (in this case always stacks)
func (d *DBWrapper) Fetch() (any, error) {
	var stacks []Stack
	result := d.db.Model(&Stack{}).Preload("Tasks").Find(&stacks)
	return stacks, result.Error
}

// Delete will remove an object
// Unscoped() is used to ensure hard delete, where stack will be removed from db instead of being just marked as "deleted"
// DB.Unscoped().Delete(&s)
func (d *DBWrapper) Delete(obj any) {
	d.db.Unscoped().Select(clause.Associations).Delete(obj)
}

// SyncJSON will sync db state to JSON
func (d *DBWrapper) SyncJSON(ctx *state.Context) {
	DumpJSONToFile(ctx)
}

// InitializeDB will setup and ready the backing store
func InitializeDB(ctx *state.Context) error {
	db, err := gorm.Open(sqlite.Open(ctx.Config.Database()), &gorm.Config{
		// Silent mode ensures that errors logs don't interfere with the view
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return err
	}
	db.AutoMigrate(&Stack{}, &Task{})

	ctx.DB = &DBWrapper{db}
	return nil
}
