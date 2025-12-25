// Package entities defines backing store objects
package entities

import (
	"github.com/enckse/mayhem/internal/state"
	"gorm.io/gorm"
)

type (
	// Entity is the core DB entity
	Entity interface {
		Save(state.Store) Entity
		Delete(state.Store)
	}
	// EntityBase is the common model for various entities
	EntityBase struct {
		gorm.Model `json:"-"`
	}
)

// EntityID gets the backing entity id
func (e EntityBase) EntityID() uint {
	return e.ID
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
