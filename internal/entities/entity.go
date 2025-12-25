// Package entities defines backing store objects
package entities

import (
	"github.com/enckse/mayhem/internal/state"
)

// Entity is the core DB entity
type Entity interface {
	Save(state.Store) Entity
	Delete(state.Store)
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
