package entities_test

import (
	"testing"

	"github.com/enckse/mayhem/internal/entities"
)

func TestEntityID(t *testing.T) {
	e := entities.Stack{}
	e.ID = 1
	if e.EntityID() != 1 {
		t.Error("invalid id")
	}
}

func TestNewStack(t *testing.T) {
	m := &mockDB{}
	e := entities.NewStack(m)
	if e.Title != "New Stack" {
		t.Error("invalid stack")
	}
}

func TestFetchStacks(t *testing.T) {
	m := &mockDB{}
	s := entities.FetchStacks(m)
	if len(s) != 1 {
		t.Error("invalid stacks")
	}
}

func TestSaveDelete(t *testing.T) {
	m := &mockDB{}
	s := entities.NewStack(m)
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
	s.Tasks = []entities.Task{{}, {}, {}, {IsFinished: true}}
	if s.OpenTasks() != 3 {
		t.Error("invalid open task count")
	}
}
