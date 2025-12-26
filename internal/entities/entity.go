// Package entities defines backing store objects
package entities

import (
	"github.com/enckse/mayhem/internal/backend"
)

type (
	// Entity is the core DB entity
	Entity interface {
		Save(backend.Store) Entity
		Delete(backend.Store)
	}
)

// FindByIndex will find an entity by an id (from a set)
func FindByIndex[T interface{ EntityID() string }](arr []T, id string) int {
	for i, val := range arr {
		if val.EntityID() == id {
			return i
		}
	}
	return -1
}
