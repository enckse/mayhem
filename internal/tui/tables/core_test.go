package tables_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/enckse/mayhem/internal/display"
	"github.com/enckse/mayhem/internal/entities"
	"github.com/enckse/mayhem/internal/state"
	"github.com/enckse/mayhem/internal/tui/tables"
)

func TestStackRows(t *testing.T) {
	s := tables.StackRows([]entities.Stack{{PendingTaskCount: 1000}, {PendingTaskCount: 10}, {Title: "empty"}})
	if fmt.Sprintf("%v", s) != "[[ [99+]] [ [ 10]] [empty      ]]" {
		t.Errorf("bad rows: %v", s)
	}
}

func TestTaskRows(t *testing.T) {
	ctx := &state.Context{}
	ctx.State.FinishedTasks = make(map[uint]bool)
	tasks := []entities.Task{{Title: "xyz", IsFinished: true}, {IsFinished: false}}
	tasks[0].ID = 0
	tasks[1].ID = 1
	s := tables.TaskRows(tasks, ctx)
	if fmt.Sprintf("%v", s) != "[[✘ xyz          -    0] [▢           -    0]]" {
		t.Errorf("bad rows: %v", s)
	}
}

func TestNew(t *testing.T) {
	s := &display.Screen{}
	res := tables.New(tables.StackColumns, display.StackTableType, s)
	if !strings.Contains(fmt.Sprintf("%v", res), "#898989") {
		t.Errorf("invalid model: %v", res)
	}
}
