package entities

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/enckse/mayhem/internal/sqlite"
	"github.com/enckse/mayhem/internal/state"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"

	// Using pure-go implementation of GORM driver to avoid CGO issues during cross-compilation
	_ "modernc.org/sqlite"
)

// DBWrapper is our sqlite backend wrapper
type DBWrapper struct {
	db        *gorm.DB
	errors    []string
	saveJSON  bool
	directory string
	maxErrors int
}

// Save will store an entity
func (d *DBWrapper) Save(obj any) {
	d.db.Save(obj).Commit()
	d.doJSON()
}

func (d *DBWrapper) doJSON() {
	if d.saveJSON {
		fileName := filepath.Join(d.directory, state.FileName+"json")
		d.log("json", DumpJSONToFile(fileName, d))
	}
}

// Create will create an entity
func (d *DBWrapper) Create(obj any) {
	d.log("create", d.db.Create(obj).Error)
	d.doJSON()
}

// Fetch will return the data (in this case always stacks)
func (d *DBWrapper) Fetch() any {
	var stacks []Stack
	result := d.db.Model(&Stack{}).Preload("Tasks").Find(&stacks)
	d.log("fetch", result.Error)
	return stacks
}

func (d *DBWrapper) log(category string, err error) {
	if err == nil {
		return
	}
	if d.maxErrors > 0 && len(d.errors) > d.maxErrors {
		d.errors = d.errors[1:]
	}
	d.errors = append(d.errors, fmt.Sprintf("[%s] %s: %v", time.Now().Format("2006-01-02T15:04:05"), category, err))
}

// Errors will get the list of errors
func (d *DBWrapper) Errors() []string {
	return d.errors
}

// Delete will remove an object
// Unscoped() is used to ensure hard delete, where stack will be removed from db instead of being just marked as "deleted"
// DB.Unscoped().Delete(&s)
func (d *DBWrapper) Delete(obj any) {
	d.db.Unscoped().Select(clause.Associations).Delete(obj).Commit()
	d.doJSON()
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

	ctx.DB = &DBWrapper{
		db:        db,
		errors:    []string{},
		saveJSON:  ctx.Config.JSON.Sync,
		directory: ctx.Config.Data.Directory,
		maxErrors: int(ctx.Config.Log.Lines),
	}
	return nil
}
