package tables_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/enckse/mayhem/internal/display"
	"github.com/enckse/mayhem/internal/entities"
	"github.com/enckse/mayhem/internal/tui/tables"
)

func TestStackRows(t *testing.T) {
	var thousandTasks, tenTasks []entities.Task
	idx := 0
	for idx <= 1000 {
		if idx < 10 {
			tenTasks = append(tenTasks, entities.Task{})
		}
		thousandTasks = append(thousandTasks, entities.Task{})
		idx++
	}
	s := tables.StackRows([]entities.Stack{{Tasks: thousandTasks}, {Tasks: tenTasks}, {Title: "empty"}})
	if fmt.Sprintf("%v", s) != "[[ [99+]] [ [ 10]] [empty      ]]" {
		t.Errorf("bad rows: %v", s)
	}
}

func TestTaskRows(t *testing.T) {
	tasks := []entities.Task{{Title: "xyz", IsFinished: true}, {IsFinished: false}}
	tasks[0].ID = 0
	tasks[1].ID = 1
	tasks[1].UpdatedAt = time.Now()
	s := tables.TaskRows(tasks, time.Time{})
	if fmt.Sprintf("%v", s) != "[[▢           -    0] [✘ xyz          -    0]]" {
		t.Errorf("bad rows: %v", s)
	}
	s = tables.TaskRows(tasks, time.Now())
	if fmt.Sprintf("%v", s) != "[[▢           -    0]]" {
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
