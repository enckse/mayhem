// Package entities defines backing store objects
package entities

import (
	// Using pure-go implementation of GORM driver to avoid CGO issues during cross-compilation
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

// Sync will save to the DB and perform any other sync operations of interest
func Sync(ctx *state.Context, obj any) {
	ctx.DB.Save(obj)
	if ctx.Config.JSON.Sync {
		ToJSON(ctx)
	}
}

// FindByIndex will find an entity by an id (from a set)
func FindByIndex[T interface{ EntityID() uint }](arr []T, id uint) int {
	for i, val := range arr {
		if val.EntityID() == id {
			return i
		}
	}
	return -1
}
