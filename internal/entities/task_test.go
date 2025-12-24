package entities_test

import (
	"testing"
	"time"

	"github.com/enckse/mayhem/internal/entities"
	"github.com/enckse/mayhem/internal/state"
)

func TestTaskEntityID(t *testing.T) {
	e := entities.Task{}
	e.ID = 1
	if e.EntityID() != 1 {
		t.Error("invalid id")
	}
}

func TestSaveDeleteTask(t *testing.T) {
	m := &mockDB{}
	ctx := &state.Context{}
	ctx.DB = m
	s := entities.Task{}
	s.Save(ctx)
	if m.last == nil {
		t.Error("no save")
	}
	s.Delete(ctx)
	if m.last != nil {
		t.Error("no delete")
	}
}

func TestSortTasks(t *testing.T) {
	s := []entities.Task{{Title: "X00", IsFinished: true, Deadline: time.Now()}, {Title: "X01", IsFinished: false, Deadline: time.Now()}, {Title: "X10", IsFinished: true}, {Title: "X11", IsFinished: false}}
	entities.SortTasks(s)
	if len(s) != 4 || s[0].Title != "X10" || s[1].Title != "X00" || s[2].Title != "X11" || s[3].Title != "X01" {
		t.Errorf("invalid sort: %v", s)
	}
}
