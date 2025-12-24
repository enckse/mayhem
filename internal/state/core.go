// Package state handles overall app state
package state

import "gorm.io/gorm"

// Context is the overall state context
type Context struct {
	DB     *gorm.DB
	Config Config
}
