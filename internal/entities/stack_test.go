package entities_test

import (
	"testing"
	"time"

	"github.com/enckse/mayhem/internal/backend"
	"github.com/enckse/mayhem/internal/entities"
)

func TestNewStack(t *testing.T) {
	m := &mockDB{}
	e := entities.NewStack(m)
	if e.Title != "New Stack" || e.ID == "" {
		t.Error("invalid stack")
	}
}

func TestFetchStacks(t *testing.T) {
	m := &mockDB{}
	s := entities.FetchStacks(m)
	if len(s) != 1 {
		t.Error("invalid stacks")
	}
	m.data = make([]backend.Data, 3)
	m.data[0].Node = entities.Stack{ID: "y", Tasks: []entities.Task{{}, {}}}
	m.data[0].Children = make(backend.Map)
	m.data[1].Node = 1
	m.data[2].Node = entities.Stack{ID: "x", Tasks: []entities.Task{{}, {}}}
	m.data[2].Children = make(backend.Map)
	m.data[2].Children["x"] = backend.Data{Node: 1}
	m.data[2].Children["y"] = backend.Data{Node: entities.Task{}}
	s = entities.FetchStacks(m)
	if len(s) != 2 {
		t.Error("invalid stacks")
	}
	if s[0].ID != "y" || s[1].ID != "x" {
		t.Error("invalid results: ids")
	}
	if len(s[0].Tasks) != 0 || len(s[1].Tasks) != 1 {
		t.Error("invalid results: tasks")
	}
}

func TestSaveDelete(t *testing.T) {
	m := &mockDB{}
	s := entities.NewStack(m)
	m = &mockDB{}
	s.Title = ""
	s.Save(m)
	if m.last != nil || m.cat != "stack" {
		t.Error("saved")
	}
	s.Title = "title"
	s.Save(m)
	if m.last == nil {
		t.Error("no save")
	}
	s.Delete(m)
	if m.last != nil {
		t.Error("no delete")
	}
}

func TestSortStacks(t *testing.T) {
	s := []entities.Stack{{Title: "X"}, {Title: "1"}}
	entities.SortStacks(s)
	if len(s) != 2 || s[0].Title != "1" || s[1].Title != "X" {
		t.Errorf("invalid sort: %v", s)
	}
}

func TestOpenTasks(t *testing.T) {
	s := entities.Stack{}
	s.Tasks = []entities.Task{{}, {}, {}, {Finished: time.Now()}}
	if s.OpenTasks() != 3 {
		t.Error("invalid open task count")
	}
}

func TestStaskEntityID(t *testing.T) {
	e := entities.Stack{}
	e.ID = "1"
	if e.EntityID() != "1" {
		t.Error("invalid id")
	}
}
