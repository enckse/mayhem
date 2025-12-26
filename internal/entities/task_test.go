package entities_test

import (
	"testing"
	"time"

	"github.com/enckse/mayhem/internal/entities"
)

func TestSaveDeleteTask(t *testing.T) {
	m := &mockDB{}
	s := entities.Task{}
	s.Save(m)
	if m.last != nil && m.cat == "task" {
		t.Error("save")
	}
	m.cat = ""
	s.Title = "x"
	s.Priority = 100
	s.Save(m)
	if m.last != nil && m.cat == "task" {
		t.Error("save")
	}
	s.Priority = 1
	s.Save(m)
	if m.last == nil {
		t.Error("no save")
	}
	s.Delete(m)
	if m.last != nil {
		t.Error("no delete")
	}
}

func TestSortTasks(t *testing.T) {
	s := []entities.Task{{Title: "X00", Finished: time.Now(), Deadline: time.Now()}, {Title: "X01", Deadline: time.Now()}, {Title: "X10", Finished: time.Now().Add(5 * time.Second)}, {Title: "X11"}}
	entities.SortTasks(s)
	if len(s) != 4 || s[0].Title != "X01" || s[1].Title != "X11" || s[2].Title != "X00" || s[3].Title != "X10" {
		t.Errorf("invalid sort: %v", s)
	}
}

func TestTaskEntityID(t *testing.T) {
	e := entities.Task{}
	e.ID = "1"
	if e.EntityID() != "1" {
		t.Error("invalid id")
	}
}
