// Package entities defines backing store objects
package entities

import (
	// Using pure-go implementation of GORM driver to avoid CGO issues during cross-compilation
	"path/filepath"

	"github.com/enckse/mayhem/internal/app"
	"github.com/enckse/mayhem/internal/state"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Entity is the core DB entity
type Entity interface {
	Save(*state.Context) Entity
	Delete(*state.Context)
}

// InitializeDB will setup and ready the backing store
func InitializeDB(ctx *state.Context) error {
	path, err := app.DataDir()
	if err != nil {
		return err
	}

	db, err := gorm.Open(sqlite.Open(filepath.Join(path, "tasks.db")), &gorm.Config{
		// Silent mode ensures that errors logs don't interfere with the view
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return err
	}
	db.AutoMigrate(&Stack{}, &Task{})

	ctx.DB = db
	return nil
}
