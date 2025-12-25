// Package entities defines backing store objects
package entities

import (
	"github.com/enckse/mayhem/internal/state"
)

// Entity is the core DB entity
type Entity interface {
	Save(*state.Context) Entity
	Delete(*state.Context)
}

// Sync will save to the DB and perform any other sync operations of interest
func Sync(ctx *state.Context, obj any) {
	ctx.DB.Save(obj)
	if ctx.Config.JSON.Sync {
		ctx.DB.SyncJSON(ctx)
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
