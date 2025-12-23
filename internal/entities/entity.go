// Package entities defines backing store objects
package entities

import (
	// Using pure-go implementation of GORM driver to avoid CGO issues during cross-compilation

	"os"
	"path/filepath"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Entity is the core DB entity
type Entity interface {
	Save() Entity
	Delete()
}

// DB is the backing database
var DB *gorm.DB

func getStorageDir() (string, error) {
	path := os.Getenv("MAYHEM_DB_PATH")
	if path != "" {
		return path, nil
	}

	const appDir = "mayhem"
	xdg := os.Getenv("XDG_CACHE_HOME")
	if xdg != "" {
		return filepath.Join(xdg, appDir), nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".cache", appDir), nil
}

// InitializeDB will setup and ready the backing store
func InitializeDB() error {
	path, err := getStorageDir()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return err
	}

	db, err := gorm.Open(sqlite.Open(filepath.Join(path, "tasks.db")), &gorm.Config{
		// Silent mode ensures that errors logs don't interfere with the view
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return err
	}
	db.AutoMigrate(&Stack{}, &Task{}, &Step{}, &RecurTask{})

	DB = db
	return nil
}
