package state_test

import (
	"testing"

	"github.com/enckse/mayhem/internal/state"
)

func TestPathExists(t *testing.T) {
	if state.PathExists("invalid") {
		t.Error("file does not exist")
	}
	if !state.PathExists("core_test.go") {
		t.Error("invalid file")
	}
}
